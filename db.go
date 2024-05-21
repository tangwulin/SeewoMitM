package main

import tiny "github.com/Yiwen-Chan/tinydb"

var db *tiny.Database

func InitDB() error {
	storage, err := RawStorage(nil)
	if err != nil {
		return err
	}

	db, err = tiny.TinyDB(storage)
	if err != nil {
		return err
	}

	return nil
}

type WriteCallback func(data any, storage *map[string]interface{}) error

type StorageRaw struct {
	storage       *map[string]interface{}
	writeCallback WriteCallback
}

func RawStorage(writeCallback WriteCallback) (*StorageRaw, error) {
	return &StorageRaw{writeCallback: writeCallback}, nil
}

func (r *StorageRaw) Read(data any) error {
	data = r.storage
	return nil
}

func (r *StorageRaw) Write(data any) error {
	*r.storage = data.(map[string]interface{})

	if r.writeCallback != nil {
		return r.writeCallback(data, r.storage)
	}
	return nil
}

func (r *StorageRaw) Close() error {
	return nil
}
