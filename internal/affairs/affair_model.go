package affairs

import (
	"SQL/internal/database"
	"fmt"
	"sync"
	"time"
)

// TransactionStatus 表示事务状态的常量
const (
	TxnPending    = iota // 事务进行中
	TxnCommitted         // 事务已提交
	TxnRolledBack        // 事务已回滚
)

type txid uint64

// Tx 表示事务的结构体
type Tx struct {
	DB           *database.XcDB
	Meta         *TxMeta
	ReadVersions map[string]uint64 // 用于存储事务读取数据的版本信息
}

type TxMeta struct {
	ID           uint64                 // 事务ID
	Status       int                    // 事务状态
	StartTime    uint64                 // 事务开始时间戳
	CommitTime   uint64                 // 事务提交时间戳
	ReadSet      map[string]interface{} // 读操作的键列表，键为键名，值为读取的数据
	WriteSet     map[string]interface{} // 写操作的键列表，键为键名，值为写入的数据
	ConflictKeys map[string]struct{}    // 冲突检测用的键集合
}

// 初始化一个事务
func (tx *Tx) init(db *database.XcDB) {
	tx.DB = db
	tx.Meta = &TxMeta{}
	tx.Meta.ID++
}

// 获取事务的id
func (tx *Tx) ID() uint64 {
	return tx.Meta.ID
}

// Oracle 表示事务管理器的结构体
type Oracle struct {
	sync.Mutex
	committedTxns   []*Tx  // 最近提交的事务列表
	globalTimestamp uint64 // 全局时间戳
}

// NewOracle 创建一个新的Oracle对象
func NewOracle() *Oracle {
	return &Oracle{
		committedTxns:   []*Tx{},
		globalTimestamp: 0,
	}
}

// BeginTransaction 开始一个新的事务
func (o *Oracle) BeginTransaction() *Tx {
	o.Lock()
	defer o.Unlock()
	o.globalTimestamp++
	return &Tx{
		Meta: &TxMeta{
			ID:           o.globalTimestamp,
			Status:       TxnPending,
			StartTime:    getCurrentTimestamp(),
			ConflictKeys: make(map[string]struct{}),
		},
	}
}

// CommitTransaction 提交事务
func (o *Oracle) CommitTransaction(txn *Tx) {
	txn.Meta.CommitTime = getCurrentTimestamp()
	txn.Meta.Status = TxnCommitted
	o.Lock()
	defer o.Unlock()
	o.committedTxns = append(o.committedTxns, txn)
}

// RollbackTransaction 回滚事务
func (o *Oracle) RollbackTransaction(txn *Tx) {
	txn.Meta.Status = TxnRolledBack
}

// AddReadKey 添加读操作的键到事务的读集合
func (t *Tx) AddReadKey(key string, data interface{}) {
	t.Meta.ReadSet[key] = data
}

// AddWriteKey 添加写操作的键到事务的写集合
func (t *Tx) AddWriteKey(key string, data interface{}) {
	t.Meta.WriteSet[key] = data
}

// AddConflictKey 添加冲突检测的键到事务的冲突键集合
func (t *Tx) AddConflictKey(key string, data interface{}) {
	t.Meta.ConflictKeys[key] = struct{}{}
}

// getCurrentTimestamp 获取当前时间戳的函数
func getCurrentTimestamp() uint64 {
	return uint64(time.Now().UnixNano())
}

// 提交事务并进行冲突检测和版本检查
func CommitAndCheckConflict(o *Oracle, txn *Tx) error {
	// 获取所有未提交的事务
	o.Lock()
	defer o.Unlock()
	uncommittedTxns := o.committedTxns

	// 检查冲突和版本
	for _, uncommittedTxn := range uncommittedTxns {
		if hasConflict(uncommittedTxn, txn) {
			return fmt.Errorf("conflict detected, cannot commit transaction")
		}
	}

	// 冲突检测通过，提交事务并更新版本信息
	o.CommitTransaction(txn)
	txn.UpdateReadVersions()
	return nil
}

// 检查两个事务是否存在冲突
func hasConflict(txn1, txn2 *Tx) bool {
	// 检查读写集合是否有交集
	for readKey := range txn1.Meta.ReadSet {
		if _, ok := txn2.Meta.WriteSet[readKey]; ok {
			return true
		}
	}

	// 检查写读冲突
	for writeKey := range txn1.Meta.WriteSet {
		if _, ok := txn2.Meta.ReadSet[writeKey]; ok {
			return true
		}
	}
	//// 检查数据版本冲突
	//for key, version := range txn1.ReadVersions {
	//	if latestVersion, ok := txn2.DB.GetVersion([]byte(key)); ok && latestVersion > version {
	//		return true
	//	}
	//}

	return false
}

// 更新事务的读取版本信息
func (txn *Tx) UpdateReadVersions() {
	//txn.ReadVersions = make(map[string]uint64)
	//for key := range txn.Meta.ReadSet {
	//	if version, ok := txn.DB.GetVersion([]byte(key)); ok {
	//		txn.ReadVersions[key] = version
	//	}
	//}
}
