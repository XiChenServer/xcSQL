package lsm

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"sync"
)

// 在进行存储的时候，内部含有数据的元信息还有就是数据的存储位置

type DataInfo struct {
	model.DataMeta
	storage.StorageLocation
}

// 跳表节点结构
type SkipListNode struct {
	Key      []byte
	DataInfo *DataInfo
	Next     []*SkipListNode // 指向下一个节点的指针数组
}
type SkipListInfo struct {
	MaxKey []byte // 这个表最大的键
	MinKey []byte // 这个表最小的键
}

// 跳表结构
type SkipList struct {
	Head         *SkipListNode // 头节点
	MaxLevel     int16         // 最大层数
	Size         uint32        // 跳表中节点数量
	mu           sync.RWMutex  // 用于保护并发访问
	SkipListInfo *SkipListInfo // 表中的一些信息
}

// 初始化跳表
func NewSkipList(maxLevel int16) *SkipList {
	head := &SkipListNode{
		Key:      nil,
		DataInfo: nil,
		Next:     make([]*SkipListNode, maxLevel),
	}
	skipListInfo := &SkipListInfo{
		MaxKey: []byte{},
		MinKey: []byte{},
	}
	return &SkipList{
		Head:         head,
		MaxLevel:     maxLevel,
		Size:         0,
		SkipListInfo: skipListInfo,
	}
}
