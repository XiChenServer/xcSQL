package alv_tree

import (
	"time"
)

// 平衡二叉树的数据节点
type nodeInfo struct {
	Key      []byte
	Type     uint16
	TTL      time.Duration
	FileName []byte
	Offset   int64
	Size     int64
}

// 平衡二叉树的整体结构
type AlvNode struct {
	data           nodeInfo
	lchild, rchild *AlvNode
}
