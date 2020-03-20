// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/cache/redis"
//   "github.com/astaxie/beego/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.me/docs/module/cache.md
package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"

	"github.com/astaxie/beego/cache"
	"strings"
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
}

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey    = "beecacheRedis"
	RedisClient   = &Cache{key: DefaultKey}
	RedisPSClient = &Cache{key: DefaultKey}
)

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

func Connect(config string) {
	RedisClient.StartAndGC(config)
}

func ConnectPS(config string) {
	RedisPSClient.StartAndGC(config)
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) doh(commandName string, key string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	key = rc.associate(key)
	c := rc.p.Get()
	defer c.Close()

	var as []interface{}
	as = append(as, key)
	for _, v := range args {
		as = append(as, v)
	}

	return c.Do(commandName, as...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetMulti get cache from redis.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SET", key, val, "EX", "300")
	return err
}

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *Cache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *Cache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// HDel delete one or more fields
func (rc *Cache) HDel(key string, fileds ...interface{}) error {
	_, err := rc.doh("HDEL", key, fileds...)
	return err
}

// HIsExist check if cached field value exists or not by key.
func (rc *Cache) HIsExist(key string, filed interface{}) bool {
	v, err := redis.Bool(rc.doh("HEXISTS", key, filed))
	if err != nil {
		return false
	}
	return v
}

// HGet get cached filed value by key.
func (rc *Cache) HGet(key string, filed interface{}) interface{} {
	if v, err := rc.doh("HGET", key, filed); err == nil {
		return v
	}
	return nil
}

// HGetAll get all cached values by key
func (rc *Cache) HGetAll(key string) []interface{} {
	values, err := redis.Values(rc.do("HGETALL", key))
	if err != nil {
		return nil
	}

	return values
}

// HMGet get all field values given by key
func (rc *Cache) HMGet(key string, fileds ...interface{}) []interface{} {
	values, err := redis.Values(rc.doh("HMGET", key, fileds...))
	if err != nil {
		return nil
	}
	return values
}

// HMSet set all field values given by key
func (rc *Cache) HMSet(key string, fvs ...interface{}) error {
	_, err := rc.doh("HMSET", key, fvs...)
	return err
}

// HSet set field value by key
func (rc *Cache) HSet(key string, filed interface{}, val interface{}) error {
	_, err := rc.doh("HSET", key, filed, val)
	return err
}

// HIncrBy atomicly incr or decr
func (rc *Cache) HIncrBy(key string, filed interface{}, delta int) error {
	_, err := rc.doh("HINCRBY", key, filed, delta)
	return err
}

// HExpire set tts
func (rc *Cache) HExpire(key string, tts int) error {
	_, err := rc.doh("EXPIRE", key, tts)
	return err
}

// Subscribe
func (rc *Cache) Subscribe(channel string, handler func(msg interface{})) {
	c := rc.p.Get()
	c.Send("SUBSCRIBE", channel)
	c.Flush()
	for {
		reply, err := c.Receive()
		if err != nil {
			logs.Error(err)
		}
		// process pushed message
		handler(reply)
	}
}

// Publish
func (rc *Cache) Publish(channel string, msg interface{}) error {
	_, err := rc.do("PUBLISH", channel, msg)
	return err
}

// RPUSH
func (rc *Cache) Rpush(key string, msg interface{}) error {
	_, err := rc.do("RPUSH", key, msg)
	return err
}

// LPOP
func (rc *Cache) Lpop(key string) (interface{}, error) {
	msg, err := rc.do("LPOP", key)
	return msg, err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
