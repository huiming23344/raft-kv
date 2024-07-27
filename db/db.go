package db

import (
	"fmt"
	"github.com/luo/kv-raft/db/cache"
	"github.com/luo/kv-raft/db/engines"
)

type DB interface {

	// Set the value of a string key to a string.
	// If the key already exists, the previous value will be overwritten.
	Set(key, value string) error

	// Get the string value of a given string key.
	// Return `None` if the given key does not exit.
	Get(key string) (string, error)

	// Remove a given key.
	//  It returns `kvserror::KeyNotFound` if the given key not found.
	Remove(key string) error
}

type db struct {
	engine engines.KvsEngine
	cache  cache.Cache
}

func NewDB(path string, cacheCap int) (DB, error) {
	engine, err := engines.NewKvsStore(path)
	if err != nil {
		return nil, err
	}
	myCache := cache.NewLRUCache(cacheCap)
	return db{
		engine: engine,
		cache:  myCache,
	}, nil
}

func (d db) Set(key, value string) error {
	if err := d.engine.Set(key, value); err != nil {
		return err
	}
	d.cache.Set(key, value)
	return nil
}

func (d db) Get(key string) (string, error) {
	if data, ok := d.cache.Get(key); ok {
		fmt.Println("get from cache")
		return data, nil
	}
	data, err := d.engine.Get(key)
	if err != nil {
		return "", err
	}
	d.cache.Set(key, data)
	return data, nil
}

func (d db) Remove(key string) error {
	if err := d.engine.Remove(key); err != nil {
		return err
	}
	fmt.Println("removed key from cache:", key)
	d.cache.Remove(key)
	return nil
}
