package database

import (
	"time"
)

// KeyValue 表示键值对
type KeyValue struct {
	DataMeta   *DataMeta
	Version    uint32    // 版本号
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 修改时间
	AccessTime time.Time // 访问时间
	DataType   uint16    // 数据类型
	DataMark   uint16    // 权限控制信息
	checksum   uint32    //校验和
}
type DataMeta struct {
	Key       []byte // 键
	Value     []byte // 值，可以根据需要选择不同的数据类型
	Extra     []byte // 其他的
	KeySize   uint32
	ValueSize uint32
	ExtraSize uint32
	TTL       time.Duration // 生存时间，0 表示永不过期
}

func NewKeyValueEntry(key, value []byte, t, m uint16, ttl ...time.Duration) *KeyValue {
	var expireTime time.Duration

	// 如果提供了ttl参数，则使用提供的ttl值
	if len(ttl) > 0 {
		expireTime = ttl[0]
	} else {
		// 否则默认设置为0，即永不过期
		expireTime = 0
	}
	return NewKeyValuePair(key, value, nil, t, m, expireTime)
}

func NewKeyValuePair(key, value, extra []byte, t, m uint16, ttl time.Duration) *KeyValue {
	// 获取当前的时间
	currentTime := time.Now()
	return &KeyValue{
		DataMeta: &DataMeta{
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
