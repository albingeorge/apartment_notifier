package main

import (
	"sync"

	"github.com/go-redis/redis/v9"
)

var lock = &sync.Mutex{}

type db struct {
	redis *redis.Client
}

var dbInstance *db

func getInstance() *db {
	if dbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		dbInstance = &db{
			redis: rdb,
		}
	}

	return dbInstance
}
