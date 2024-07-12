# kvs设计

## 




### 数据的存储方式
command存储记录，command存储记录的位置
```go
type CommandType string

type Command struct {
Type  CommandType `json:"type"`
Key   string      `json:"key"`
Value string      `json:"value,omitempty"`
}

// command position used to find command in logFiles
type CommandPos struct {
gen uint64
pos uint64
len uint64
}

func NewCommandPos(gen, start, end uint64) *CommandPos {
return &CommandPos{gen, start, end - start}
}
```


### 读写文件
每一个read/write buffer with pos 是一个文件流，存储在文件中
```go
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
```

### kvs-set


```go
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
```