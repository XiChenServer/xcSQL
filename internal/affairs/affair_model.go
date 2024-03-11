package affairs

import (
	"sync"
	"time"
)

// TransactionStatus 表示事务状态的常量
const (
	TxnPending    = iota // 事务进行中
	TxnCommitted         // 事务已提交
	TxnRolledBack        // 事务已回滚
)

// Transaction 表示事务的结构体
type Transaction struct {
	ID           uint64              // 事务ID
	Status       int                 // 事务状态
	StartTime    uint64              // 事务开始时间戳
	CommitTime   uint64              // 事务提交时间戳
	ReadSet      [][]byte            // 读操作的键列表
	WriteSet     [][]byte            // 写操作的键列表
	ConflictKeys map[string]struct{} // 冲突检测用的键集合
}

// Oracle 表示事务管理器的结构体
type Oracle struct {
	sync.Mutex
	committedTxns   []*Transaction // 最近提交的事务列表
	globalTimestamp uint64         // 全局时间戳
}

// NewOracle 创建一个新的Oracle对象
func NewOracle() *Oracle {
	return &Oracle{
		committedTxns:   []*Transaction{},
		globalTimestamp: 0,
	}
}

// BeginTransaction 开始一个新的事务
func (o *Oracle) BeginTransaction() *Transaction {
	o.Lock()
	defer o.Unlock()
	o.globalTimestamp++
	return &Transaction{
		ID:           o.globalTimestamp,
		Status:       TxnPending,
		StartTime:    getCurrentTimestamp(),
		ConflictKeys: make(map[string]struct{}),
	}
}

// CommitTransaction 提交事务
func (o *Oracle) CommitTransaction(txn *Transaction) {
	txn.CommitTime = getCurrentTimestamp()
	txn.Status = TxnCommitted
	o.Lock()
	defer o.Unlock()
	o.committedTxns = append(o.committedTxns, txn)
}

// RollbackTransaction 回滚事务
func (o *Oracle) RollbackTransaction(txn *Transaction) {
	txn.Status = TxnRolledBack
}

// AddReadKey 添加读操作的键到事务的读集合
func (t *Transaction) AddReadKey(key []byte) {
	t.ReadSet = append(t.ReadSet, key)
}

// AddWriteKey 添加写操作的键到事务的写集合
func (t *Transaction) AddWriteKey(key []byte) {
	t.WriteSet = append(t.WriteSet, key)
}

// AddConflictKey 添加冲突检测的键到事务的冲突键集合
func (t *Transaction) AddConflictKey(key []byte) {
	t.ConflictKeys[string(key)] = struct{}{}
}

// getCurrentTimestamp 获取当前时间戳的函数
func getCurrentTimestamp() uint64 {
	return uint64(time.Now().UnixNano())
}
