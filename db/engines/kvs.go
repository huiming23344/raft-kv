package engines

import (
	"bufio"
	"encoding/json"
	"fmt"
	errs "github.com/luo/kv-raft/errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const CompactionThreshold uint64 = 1024 * 1024

type CommandType string

const (
	SET    CommandType = "SET"
	DELETE CommandType = "DELETE"
)

type Command struct {
	Type  CommandType `json:"type"`
	Key   string      `json:"key"`
	Value string      `json:"value,omitempty"`
}

// CommandPos is a position used to find command in logFiles
type CommandPos struct {
	gen uint64
	pos uint64
	len uint64
}

func NewCommandPos(gen, start, end uint64) *CommandPos {
	return &CommandPos{gen, start, end - start}
}

type BufReaderWithPos struct {
	file *os.File
	pos  uint64
}

func NewBufReaderWithPos(file *os.File) *BufReaderWithPos {
	reader := &BufReaderWithPos{
		file: file,
		pos:  0,
	}
	return reader
}

func (r *BufReaderWithPos) readCommand(commandPos *CommandPos) (*Command, error) {
	buf := make([]byte, commandPos.len)
	err := r.read(buf)
	if err != nil {
		return nil, err
	}
	cmd := &Command{}
	err = json.Unmarshal(buf, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (r *BufReaderWithPos) read(buf []byte) error {
	n, err := r.file.Read(buf)
	if err != nil {
		return err
	}
	r.pos += uint64(n)
	return nil
}

func (r *BufReaderWithPos) seek(pos uint64) error {
	n, err := r.file.Seek(int64(pos), io.SeekStart)
	r.pos = uint64(n)
	return err
}

type BufWriterWithPos struct {
	buf *bufio.Writer
	pos uint64
}

func (w *BufWriterWithPos) write(p []byte) error {
	n, err := w.buf.Write(p)
	if err != nil {
		return err
	}
	w.pos += uint64(n)
	return nil
}

func (w *BufWriterWithPos) flush() error {
	return w.buf.Flush()
}

func NewBufWriterWithPos(file *os.File) *BufWriterWithPos {
	return &BufWriterWithPos{buf: bufio.NewWriter(file), pos: 0}
}

type KvsStore struct {
	mutex       sync.Mutex
	path        string
	readers     map[uint64]*BufReaderWithPos
	writer      *BufWriterWithPos
	currentGen  uint64
	index       *sync.Map
	unCompacted uint64 // record useless bytes
}

func NewKvsStore(path string) (KvsEngine, error) {
	genList, err := sortedGenList(path)
	if err != nil {
		return nil, err
	}
	index := &sync.Map{}
	readers := make(map[uint64]*BufReaderWithPos)
	var uncompacted uint64
	for _, gen := range genList {
		file, err := os.Open(logPath(path, uint64(gen)))
		if err != nil {
			return nil, err
		}
		reader := NewBufReaderWithPos(file)
		n, err := load(uint64(gen), reader, index)
		if err != nil {
			return nil, err
		}
		uncompacted += n
		readers[uint64(gen)] = reader
	}
	currentGen := 1
	if len(genList) > 0 {
		currentGen = genList[len(genList)-1] + 1
	}
	writer, err := newLogFile(path, uint64(currentGen))
	if err != nil {
		return nil, err
	}
	kvsStore := &KvsStore{
		path:        path,
		readers:     readers,
		writer:      writer,
		currentGen:  uint64(currentGen),
		index:       index,
		unCompacted: uncompacted,
	}
	return kvsStore, nil
}

func load(gen uint64, reader *BufReaderWithPos, index *sync.Map) (uint64, error) {
	var uncompacted, pos uint64
	buf := bufio.NewReader(reader.file)
	for bytes, err := buf.ReadSlice('}'); err == nil; bytes, err = buf.ReadSlice('}') {
		newPos := pos + uint64(len(bytes))
		cmd := &Command{}
		err := json.Unmarshal(bytes, cmd)
		if err != nil {
			return 0, err
		}
		if cmd.Type == SET {
			if val, ok := index.Load(cmd.Key); ok {
				uncompacted += val.(*CommandPos).len
			}
			cmdPos := NewCommandPos(gen, pos, newPos)
			index.Store(cmd.Key, cmdPos)
		}
		if cmd.Type == DELETE {
			if val, ok := index.Load(cmd.Key); ok {
				uncompacted += val.(*CommandPos).len
				index.Delete(cmd.Key)
			}
			// Remove命令在下一次压缩中删除，因此将长度置为未压缩
			uncompacted += newPos - pos
		}
		pos = newPos
	}
	return uncompacted, nil
}

func logPath(path string, gen uint64) string {
	return fmt.Sprintf("%s/%d.log", path, gen)
}

func newLogFile(path string, gen uint64) (*BufWriterWithPos, error) {
	file, err := os.Create(logPath(path, gen))
	if err != nil {
		return nil, err
	}
	writer := NewBufWriterWithPos(file)
	return writer, nil
}

func sortedGenList(path string) ([]int, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}
	genList := make([]int, 0)
	files, err := os.ReadDir(path)
	if err != nil {
		return genList, err
	}
	for _, v := range files {
		if strings.HasSuffix(v.Name(), ".log") {
			tempArr := strings.Split(v.Name(), ".")
			if len(tempArr) != 2 {
				continue
			}
			seq, err := strconv.ParseUint(tempArr[0], 10, 31)
			if err != nil {
				continue
			}
			genList = append(genList, int(seq))
		}
	}
	sort.Ints(genList)
	return genList, nil
}

// compact logFiles which will
// replace all the stale logs with compactionGen
// and update kvs.currentGen += 2
func (kvs *KvsStore) compact() (err error) {
	compactionGen := kvs.currentGen + 1
	kvs.currentGen += 2
	kvs.writer, err = newLogFile(kvs.path, kvs.currentGen)
	if err != nil {
		return err
	}
	compactionWriter, err := newLogFile(kvs.path, compactionGen)
	if err != nil {
		return err
	}
	kvs.index.Range(func(key, value interface{}) bool {
		cmdPos := value.(*CommandPos)
		reader := kvs.readers[cmdPos.gen]
		if reader.pos != cmdPos.pos {
			if err = reader.seek(cmdPos.pos); err != nil {
				return false
			}
		}
		buf := make([]byte, cmdPos.len)
		if err = reader.read(buf); err != nil {
			return false
		}
		if err = compactionWriter.write(buf); err != nil {
			return false
		}
		return true
	})
	if err = compactionWriter.flush(); err != nil {
		return err
	}

	// remove stale log files
	genList, err := sortedGenList(kvs.path)
	if err != nil {
		return err
	}
	for _, gen := range genList {
		if uint64(gen) < compactionGen {
			err := os.Remove(logPath(kvs.path, uint64(gen)))
			if err != nil {
				return err
			}
		}
	}
	kvs.unCompacted = 0
	return nil
}

// Set will write new kv-set to current kvs.writer
// and then update kvs.index with CommandPos
// if index has older version update kvs.unCompacted
func (kvs *KvsStore) Set(key, value string) error {
	kvs.mutex.Lock()
	defer kvs.mutex.Unlock()
	cmd := &Command{SET, key, value}
	pos := kvs.writer.pos
	bytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	if err = kvs.writer.write(bytes); err != nil {
		return err
	}
	if err = kvs.writer.flush(); err != nil {
		return err
	}
	if val, ok := kvs.index.Load(key); ok {
		kvs.unCompacted += val.(*CommandPos).len
	}
	// always update the index with new kvSet
	commandPos := NewCommandPos(kvs.currentGen, pos, kvs.writer.pos)
	kvs.index.Store(key, commandPos)
	if kvs.unCompacted > CompactionThreshold {
		_ = kvs.compact()
	}
	return nil
}

func (kvs *KvsStore) Remove(key string) error {
	kvs.mutex.Lock()
	defer kvs.mutex.Unlock()
	if _, ok := kvs.index.Load(key); ok {
		cmd := &Command{DELETE, key, ""}
		pos := kvs.writer.pos
		bytes, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		err = kvs.writer.write(bytes)
		if err != nil {
			return err
		}
		err = kvs.writer.flush()
		if err != nil {
			return err
		}
		kvs.index.Delete(key)
		// Remove命令在下一次压缩中删除，因此将长度置为未压缩
		kvs.unCompacted += kvs.writer.pos - pos
		if kvs.unCompacted > CompactionThreshold {
			_ = kvs.compact()
		}
	} else {
		return errs.KeyNotFound
	}
	return nil
}

func (kvs *KvsStore) Get(key string) (string, error) {
	kvs.mutex.Lock()
	defer kvs.mutex.Unlock()

	if val, ok := kvs.index.Load(key); ok {
		pos := val.(*CommandPos)
		reader := kvs.readers[pos.gen]

		if reader == nil {
			file, err := os.Open(logPath(kvs.path, pos.gen))
			if err != nil {
				return "", err
			}
			reader = NewBufReaderWithPos(file)
			kvs.readers[pos.gen] = reader
		}

		err := reader.seek(pos.pos)
		if err != nil {
			return "", err
		}
		cmd, err := reader.readCommand(pos)
		if err != nil {
			return "", err
		}
		return cmd.Value, nil
	} else {
		return "", errs.KeyNotFound
	}
}
