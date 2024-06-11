package affairs

import (
	"SQL/internal/database"
	"SQL/internal/model"
	"fmt"
	"strconv"
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
	Snapshot     database.XcDB          //生成一个快照
}

// ID 获取事务的id
func (tx *Tx) ID() uint64 {
	return tx.Meta.ID
}

// Oracle 表示事务管理器的结构体
type Oracle struct {
	sync.Mutex
	committedTxns   []*Tx  // 最近提交的事务列表
	globalTimestamp uint64 // 全局时间戳
	db              *database.XcDB
}

// NewOracle 创建一个新的Oracle对象
func NewOracle() *Oracle {
	return &Oracle{
		committedTxns:   []*Tx{},
		globalTimestamp: 0,
	}
}

type dataTx struct {
	Key         []byte //key键
	Version     uint32 // 版本号
	OperateTime time.Time
	DataType    uint16      // 数据类型
	DataMark    uint16      // 权限控制信息
	Value       interface{} // 值，可以根据需要选择不同的数据类型
	TTL         uint64      // 生存时间，0 表示永不过期
}

// BeginTransaction 开始一个新的事务
func (o *Oracle) BeginTransaction() *Tx {
	o.Lock()
	defer o.Unlock()
	o.globalTimestamp++
	tx := &Tx{
		Meta: &TxMeta{
			ID:           o.globalTimestamp,
			Status:       TxnPending,
			StartTime:    getCurrentTimestamp(),
			ConflictKeys: make(map[string]struct{}),
			Snapshot:     *o.db,
		},
	}
	o.committedTxns = append(o.committedTxns, tx)
	o.db.Wal.Write(tx) //开启事务的时候也是写入
	return tx
}

// CommitTransaction 提交事务
func (o *Oracle) CommitTransaction(txn *Tx) {
	txn.Meta.CommitTime = getCurrentTimestamp()
	txn.Meta.Status = TxnCommitted
	o.Lock()
	defer o.Unlock()
	o.committedTxns = append(o.committedTxns, txn)
	// 记录事务提交到 Redo Log
	o.db.Wal.Write(txn)
}

// RollbackTransaction 回滚事务
func (o *Oracle) RollbackTransaction(txn *Tx) {
	txn.Meta.Status = TxnRolledBack
	o.db.Wal.Write(txn) //开启事务的时候也是写入
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

// CommitAndCheckConflict 提交事务并进行冲突检测和版本检查
func CommitAndCheckConflict(o *Oracle, txn *Tx) error {
	// 获取所有未提交的事务
	o.Lock()
	defer o.Unlock()
	uncommittedTxns := o.committedTxns

	// 检查冲突和版本
	for _, uncommittedTxn := range uncommittedTxns {
		if hasConflict(uncommittedTxn, txn) {
			// 如果检测到冲突，则回滚事务并记录到 Redo Log
			o.RollbackTransaction(txn)
			return fmt.Errorf("conflict detected, transaction rolled back")
		}
	}

	// 冲突检测通过，提交事务并更新版本信息
	o.CommitTransaction(txn)
	txn.UpdateReadVersions()

	return nil
}

// hasConflict 检查两个事务是否存在冲突
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

	return false
}

// UpdateReadVersions 更新事务的读取版本信息
func (txn *Tx) UpdateReadVersions() {
	for key, v := range txn.Meta.WriteSet {
		data := v.(dataTx)
		switch data.DataType {
		case model.XCDB_List:
			switch data.DataMark {
			case model.XCDB_ListLPOP:
				txn.DB.LPOP(data.Key)
			case model.XCDB_ListLPUSH:
				txn.DB.LPUSH(data.Key, data.Value.([][]byte), data.TTL)
			case model.XCDB_ListRPOP:
				txn.DB.RPOP(data.Key)
			case model.XCDB_RPUSH:
				txn.DB.RPUSH(data.Key, data.Value.([][]byte), data.TTL)
			}
		case model.XCDB_String:
			switch data.DataMark {
			case model.XCDB_StringSet:
				txn.DB.Set([]byte(key), data.Value.([]byte), data.TTL)
			case model.XCDB_Append:
				txn.DB.Append([]byte(key), data.Value.([]byte))
			}
		case model.XCDB_Hash:
			switch data.DataMark {
			case model.XCDB_HSet:
				txn.DB.HSet(data.Key, data.Value.(map[string]string))
			case model.XCDB_HDel:
				txn.DB.HDel(data.Key, data.Value.([]string)...)
			}
		case model.XCDB_Set:
			switch data.DataMark {
			case model.XCDB_SetSADD:
				txn.DB.SAdd(data.Key, data.Value.([][]byte))
			case model.XCDB_SetSREM:
				txn.DB.SRem(data.Key, data.Value.([][]byte))
			}
		}
	}
}

// TxGet 事务的字符串读取操作
func (txn *Tx) TxGet(key []byte) (interface{}, error) {
	num := strconv.Itoa(int(model.XCDB_StringGet))
	// 先从事务的 WriteSet 中查找是否有对应的数据
	if val, ok := txn.Meta.WriteSet[num+string(key)]; ok {
		return []byte(val.(string)), nil
	}

	// 再从事务的 ReadSet 中查找是否有对应的数据
	if val, ok := txn.Meta.ReadSet[num+string(key)]; ok {
		return []byte(val.(string)), nil
	}

	// 最后到快照里面获取
	data, err := txn.Meta.Snapshot.Get(key)
	if err != nil {
		return nil, err
	}

	// 更新事务的 ReadSet
	DataTx := dataTx{
		Key:         key,
		Value:       []byte(data),
		OperateTime: time.Now(),
		DataType:    model.XCDB_StringGet,
		DataMark:    model.XCDB_String,
		TTL:         0,
	}
	txn.Meta.ReadSet[num+string(key)] = DataTx
	txn.DB.Wal.Write(DataTx) //将数据写入到redolog之中
	return data, nil
}
