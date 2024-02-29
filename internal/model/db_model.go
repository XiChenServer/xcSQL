package model

import (
	"time"
)

const (
	List      uint16 = 11
	ListLPUSH uint16 = 12
	String    uint16 = 1
	StringSet uint16 = 2
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
