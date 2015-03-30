// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package utils

import (
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

var defaultTimeout = 1 * time.Second

// get redis's slot size
func SlotsInfo(addr string, fromSlot, toSlot int) (map[int]int, error) {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	var reply []interface{}
	var val []interface{}

	reply, err = redis.Values(c.Do("SLOTSINFO", fromSlot, toSlot-fromSlot+1))
	if err != nil {
		return nil, err
	}

	ret := make(map[int]int)
	for {
		if reply == nil || len(reply) == 0 {
			break
		}
		if reply, err = redis.Scan(reply, &val); err != nil {
			return nil, err
		}
		var slot, keyCount int
		_, err := redis.Scan(val, &slot, &keyCount)
		if err != nil {
			return nil, err
		}
		ret[slot] = keyCount
	}
	return ret, nil
}

func GetRedisStat(addr string) (map[string]string, error) {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	ret, err := redis.String(c.Do("INFO"))
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	lines := strings.Split(ret, "\n")
	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			m[k] = v
		}
	}

	var reply []string

	reply, err = redis.Strings(c.Do("config", "get", "maxmemory"))
	if err != nil {
		return nil, err
	}
	// we got result
	if len(reply) == 2 {
		if reply[1] != "0" {
			m["maxmemory"] = reply[1]
		} else {
			m["maxmemory"] = "∞"
		}
	}

	return m, nil
}

func GetRedisConfig(addr string, configName string) (string, error) {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return "", err
	}
	defer c.Close()
	ret, err := redis.Strings(c.Do("config", "get", configName))
	if err != nil {
		return "", err
	}
	if len(ret) == 2 {
		return ret[1], nil
	}
	return "", nil
}

func SlaveNoOne(addr string) error {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Do("SLAVEOF", "NO", "ONE")
	if err != nil {
		return err
	}
	return nil
}

func OpAof(addr string, on bool) error {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return err
	}
	defer c.Close()

	if on {
		_, err = c.Do("config", "set", "appendonly", "yes")
	} else {
		_, err = c.Do("config", "set", "appendonly", "no")
	}
	if err != nil {
		return err
	}
	return nil
}

func SetSlaveOf(addrMaster string, addrSlave string) error {
	c, err := redis.DialTimeout("tcp", addrSlave, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return err
	}
	defer c.Close()
	temp := strings.Split(addrMaster, ":")
	ip := temp[0]
	port := temp[1]
	_, err = c.Do("SLAVEOF", ip, port)
	if err != nil {
		return err
	}
	return nil
}

func Ping(addr string) error {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return err
	}
	defer c.Close()
	return err
}

func CloseRdb(addr string) error {
	c, err := redis.DialTimeout("tcp", addr, defaultTimeout, defaultTimeout, defaultTimeout)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("config", "set", "save", "")
	if err != nil {
		return err
	}
	return nil
}
