package lsm

import (
	"SQL/internal/database"
	"SQL/internal/storage"
	"bufio"
	"bytes"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

// 在进行存储的时候，内部含有数据的元信息还有就是数据的存储位置

type DataInfo struct {
	database.DataMeta
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

// 新建一个跳表
func NewSkipListLevel() *SkipList {
	return NewSkipList(16) // 初始化一个最大层数为 3 的跳表
}

// 初始化跳表
func NewSkipList(maxLevel int16) *SkipList {
	head := &SkipListNode{
		Key:      nil,
		DataInfo: nil,
		Next:     make([]*SkipListNode, maxLevel),
	}
	skipListInfo := &SkipListInfo{
		MaxKey: nil,
		MinKey: nil,
	}
	return &SkipList{
		Head:         head,
		MaxLevel:     maxLevel,
		Size:         0,
		SkipListInfo: skipListInfo,
	}
}

// 在跳表中按照字典顺序插入节点
func (sl *SkipList) InsertInOrder(key []byte, value *DataInfo) {

	// 检查待插入的键是否已经存在
	existingNode := sl.Search(key)
	if existingNode != nil {
		// 如果键已经存在，更新相应的值
		existingNode.DataInfo = value
		return
	}

	// 创建新的节点
	newNode := &SkipListNode{
		Key:      key,
		DataInfo: value,
		Next:     make([]*SkipListNode, sl.randomLevel()+1), // 增加一级以保证跳表的平衡性
	}

	// 获取更新路径
	update := make([]*SkipListNode, len(newNode.Next))
	node := sl.Head
	for i := len(newNode.Next) - 1; i >= 0; i-- {
		for node.Next[i] != nil && bytes.Compare(node.Next[i].Key, key) < 0 {
			node = node.Next[i]
		}
		update[i] = node
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	// 更新最大键和最小键
	if sl.Size == 0 {
		sl.SkipListInfo.MaxKey = key
		sl.SkipListInfo.MinKey = key
	}
	if bytes.Compare(key, sl.SkipListInfo.MaxKey) > 0 {
		sl.SkipListInfo.MaxKey = key
	}
	if bytes.Compare(key, sl.SkipListInfo.MinKey) < 0 {
		sl.SkipListInfo.MinKey = key
	}
	// 插入新节点
	for i := 0; i < len(newNode.Next); i++ {
		// 如果下一个节点存在且键大于待插入键，将新节点插入到当前节点之后
		if update[i].Next[i] != nil && bytes.Compare(update[i].Next[i].Key, key) > 0 {
			newNode.Next[i] = update[i].Next[i]
			update[i].Next[i] = newNode
		} else {
			// 否则，将新节点插入到当前节点之前
			newNode.Next[i] = update[i].Next[i]
			update[i].Next[i] = newNode
		}
	}

	// 增加跳表的大小
	sl.Size++
}

// 在跳表中查找节点
func (sl *SkipList) Search(key []byte) *SkipListNode {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	node := sl.Head
	for i := sl.MaxLevel - 1; i >= 0; i-- {
		for node.Next[i] != nil && string(node.Next[i].Key) < string(key) {
			node = node.Next[i]
		}
	}
	if node.Next[0] != nil && string(node.Next[0].Key) == string(key) {
		return node.Next[0]
	}
	return nil
}

// 随机生成节点层数
func (sl *SkipList) randomLevel() int {
	level := 1
	for level < int(sl.MaxLevel) && rand.Float32() < 0.5 {
		level++
	}
	return level
}

// 将跳表中的所有数据打印并保存到文件
func (sl *SkipList) PrintToFile(filePath string) error {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	// 打开文件准备写入
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 bufio.Writer 提高写入性能
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 遍历跳表中的所有节点并写入文件
	node := sl.Head.Next[0] // 跳过头节点
	for node != nil {
		line := "Key: " + string(node.Key) + ", Value: " + string(node.DataInfo.Value) + ", Extra: " + string(node.DataInfo.Extra) + ", TTL: " + strconv.FormatInt(int64(node.DataInfo.TTL.Seconds()), 10) + "\n"
		writer.WriteString(line)
		node = node.Next[0]
	}

	return nil
}

// 遍历跳表中的每个节点，并对每个节点执行指定的操作
func (sl *SkipList) ForEach(f func(key []byte, value *DataInfo) bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	node := sl.Head.Next[0]
	for node != nil {
		if !f(node.Key, node.DataInfo) {
			break
		}
		node = node.Next[0]
	}
}
