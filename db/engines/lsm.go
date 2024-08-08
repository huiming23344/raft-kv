package engines

import (
	"errors"
	"github.com/huiming23344/kv-raft/config"
	"github.com/huiming23344/kv-raft/db/engines/lsm"
	lsmcfg "github.com/huiming23344/kv-raft/db/engines/lsm/config"
)

type lsmEngine struct {
}

var _ KvsEngine = (*lsmEngine)(nil)

func NewLsmEngine(path string) KvsEngine {
	InitLsmEngine(path)
	return &lsmEngine{}
}

func InitLsmEngine(path string) {
	cfg := config.GlobalConfig()
	lsm.Start(lsmcfg.Config{
		DataDir:          path,
		Level0Size:       cfg.Lsm.Level0Size,
		PartSize:         cfg.Lsm.PartSize,
		Threshold:        cfg.Lsm.Threshold,
		CheckInterval:    cfg.Lsm.CheckInterval,
		CompressInterval: cfg.Lsm.CompressInterval,
	})
}

func (l *lsmEngine) Set(key, value string) error {
	if success := lsm.Set[string](key, value); !success {
		return errors.New("set failed")
	}
	return nil
}

func (l *lsmEngine) Remove(key string) error {
	lsm.Delete[string](key)
	return nil
}

func (l *lsmEngine) Get(key string) (string, error) {
	value, success := lsm.Get[string](key)
	if !success {
		return "", errors.New("key not found")
	}
	return value, nil
}
