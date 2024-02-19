package database

import (
	"sync"
	"time"
)

// KeyValuePair 表示键值对
type KeyValuePair struct {
	KvType          *KvType
	TTL             time.Duration // 生存时间，0 表示永不过期
	Version         uint32        // 版本号
	CreateTime      time.Time     // 创建时间
	UpdateTime      time.Time     // 修改时间
	AccessTime      time.Time     // 访问时间
	Tags            []byte        // 标签
	DataType        uint16        // 数据类型
	Permission      uint16        // 权限控制信息
	StorageLocation uint16        // 存储位置
	// 读写锁，用于并发读写控制
	sync.RWMutex
}
type KvType struct {
	Key       []byte // 键
	Value     []byte // 值，可以根据需要选择不同的数据类型
	Extra     []byte // 其他的
	KeySize   uint32
	ValueSize uint32
	ExtraSize uint32
}
