package model

import (
	"time"
)

const (
	XCDB_String    uint16 = 1
	XCDB_StringSet uint16 = 2
	XCDB_StringGet uint16 = 3
	XCDB_Strlen    uint16 = 4
	XCDB_Append    uint16 = 5
)
const (
	XCDB_List      uint16 = 11
	XCDB_ListLPUSH uint16 = 12
	XCDB_ListLPOP  uint16 = 13
	XCDB_ListRPOP  uint16 = 14
	XCDB_RPUSH     uint16 = 15
	XCDB_LRANGE    uint16 = 16
	XCDB_LINDEX    uint16 = 17
	XCDB_LLEN      uint16 = 18
)
const (
	XCDB_Hash     uint16 = 21
	XCDB_HSet     uint16 = 22
	XCDB_HGet     uint16 = 23
	XCDB_HGETALL  uint16 = 24
	XCDB_HDel     uint16 = 25
	XCDB_HExistst uint16 = 26
	XCDB_HKeys    uint16 = 27
	XCDB_HVals    uint16 = 28
	XCDB_HLen     uint16 = 29
)

const (
	XCDB_Set       = 31
	XCDB_SetSADD   = 32
	XCDB_SetSREM   = 33
	XCDB_SMembers  = 34
	XCDB_SIsMember = 35
	XCDB_SCard     = 36
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
	Value      []byte    // 值，可以根据需要选择不同的数据类型
	ValueSize  uint32
}
type DataMeta struct {
	Key       []byte // 键
	Extra     []byte // 其他的
	KeySize   uint32
	ExtraSize uint32
	TTL       time.Duration // 生存时间，0 表示永不过期
}
