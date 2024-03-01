package database

import (
	"SQL/internal/model"
	"time"
)

func NewKeyValueEntry(key, value []byte, t, m uint16, ttl ...time.Duration) *model.KeyValue {
	var expireTime time.Duration

	// 如果提供了ttl参数，则使用提供的ttl值
	if len(ttl) > 0 {
		expireTime = ttl[0]
	} else {
		// 否则默认设置为0，即永不过期
		expireTime = 0
	}
	var extra []byte
	if t == model.XCDB_List {
		extra = []byte("XCDB_String")
	} else if t == model.XCDB_List {
		extra = []byte("XCDB_List")
	}
	return NewKeyValuePair(key, value, extra, t, m, expireTime)
}

func NewKeyValuePair(key, value, extra []byte, t, m uint16, ttl time.Duration) *model.KeyValue {
	// 获取当前的时间
	currentTime := time.Now()
	return &model.KeyValue{
		DataMeta: &model.DataMeta{
			Key:       key,
			Value:     value,
			Extra:     extra,
			KeySize:   uint32(len(key)),
			ValueSize: uint32(len(key)),
			ExtraSize: uint32(len(key)),
			TTL:       ttl,
		},
		Version:    0,
		CreateTime: currentTime,
		UpdateTime: currentTime,
		AccessTime: currentTime,
		DataType:   t,
		DataMark:   m,
	}
}
