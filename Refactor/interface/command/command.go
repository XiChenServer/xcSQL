package command

import (
	"fmt"
	"time"
)

type Command interface {
	Set(string, string, time.Duration) error
	Del(string) error
	Get(string) (string, error)
}

type StringCommand struct {
	store map[string]kvEntry
}

type kvEntry struct {
	value      string
	expiration time.Time
}

func NewStringCommand() *StringCommand {
	return &StringCommand{
		store: make(map[string]kvEntry),
	}
}

func (c *StringCommand) Set(key string, value string, expireTime time.Duration) error {
	expiration := time.Time{}
	if expireTime > 0 {
		expiration = time.Now().Add(expireTime)
	}
	c.store[key] = kvEntry{value: value, expiration: expiration}
	return nil
}

func (c *StringCommand) Del(key string) error {
	if _, exists := c.store[key]; !exists {
		return fmt.Errorf("键 %s 不存在", key)
	}
	delete(c.store, key)
	return nil
}

func (c *StringCommand) Get(key string) (string, error) {
	entry, exists := c.store[key]
	if !exists {
		return "", fmt.Errorf("键 %s 不存在", key)
	}
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		delete(c.store, key)
		return "", fmt.Errorf("键 %s 已过期", key)
	}
	return entry.value, nil
}
