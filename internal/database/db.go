package database

import (
	"SQL/internal/log"
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/internal/wal"
	"SQL/logs"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sync"
)

type XcDB struct {
	StorageManager *storage.StorageManager
	Lsm            *map[uint16]*lsm.LSMTree //lsm树
	Wal            *wal.WAL                 // redo.log之类的，有关事务的
	BinLog         *log.BinlogFile
	// 读写锁，用于并发读写控制
	Mu sync.RWMutex
	//TX   *affairs.Oracle //事务的操作
}

func NewXcDB(name string) (*XcDB, error) {
	var lsmMap = make(map[uint16]*lsm.LSMTree, 4)
	// 启动一个协程来初始化字符串类型的LSM树
	//go func() {
	//	lsmString := lsm.NewLSMTree(16, 10000, model.XCDB_String)
	//	// 在这里可以对 lsmString 进行操作，例如插入初始数据等
	//	lsmMap[model.XCDB_String] = *lsmString
	//}()
	//
	//// 启动一个协程来初始化列表类型的LSM树
	//go func() {
	//	lsmList := lsm.NewLSMTree(16, 10000, model.XCDB_List)
	//	// 在这里可以对 lsmList 进行操作，例如插入初始数据等
	//	lsmMap[model.XCDB_List] = *lsmList
	//}()

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 初始化 Viper
	v := viper.New()

	// 设置配置文件名
	v.SetConfigFile(path + "/../../config/config.yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return nil, err
	}

	maxActiveSize := v.GetUint32("lsmTree.maxActiveSize")
	maxDiskTableSize := v.GetUint32("lsmTree.maxDiskTableSize")
	//fmt.Println(maxActiveSize, maxDiskTableSize)

	lsmString := lsm.NewLSMTree(maxActiveSize, maxDiskTableSize, model.XCDB_String, name)
	lsmList := lsm.NewLSMTree(maxActiveSize, maxDiskTableSize, model.XCDB_List, name)
	lsmHash := lsm.NewLSMTree(maxActiveSize, maxDiskTableSize, model.XCDB_Hash, name)
	lsmSet := lsm.NewLSMTree(maxActiveSize, maxDiskTableSize, model.XCDB_Set, name)
	lsmMap[model.XCDB_List] = lsmList
	lsmMap[model.XCDB_String] = lsmString
	lsmMap[model.XCDB_Hash] = lsmHash
	lsmMap[model.XCDB_Set] = lsmSet

	storageManager, err := storage.LoadStorageManager(("../../data/testdata/manager/") + name + ("/disk/config.txt"))
	if err != nil {

		storageManager, err = storage.NewStorageManager("../../data/testdata/manager/"+name+"/disk", 10*1024) // 1MB 文件大小限制
		if err != nil {
			logs.SugarLogger.Error("failed to create storage manager: %v", err)
			return nil, err
		}

	}

	wal, err := wal.NewWAL("../../data/testdata/manager/"+name+"/wal.log", "../../data/testdata/manager/"+name+"/walInfo.log")

	binlog, err := log.NewBinlogFile(name)
	//txOracle := affairs.NewOracle()
	if err != nil {
		logs.SugarLogger.Error("wal.log create fail")
		return nil, err
	}
	return &XcDB{
		Lsm:            &lsmMap,
		StorageManager: storageManager,
		Mu:             sync.RWMutex{},
		Wal:            wal,
		BinLog:         binlog,
		//TX:             txOracle,
	}, nil
}

func DBConnect(name string) *XcDB {
	// 初始化日志记录器
	logs.InitLogger(name)

	// 连接数据库
	db, err := NewXcDB(name)
	if err != nil {
		logs.SugarLogger.Panic("new db fail")
	}

	return db
}
func DBExit(db *XcDB) error {
	// 在退出时保存活动数据到磁盘并将磁盘数据打印到文件中以供 LSM 树使用
	// 将存储管理器配置保存到文件
	err := storage.SaveStorageManager(db.StorageManager, string(db.StorageManager.StoragePath)+"/config.txt")
	//err := storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
	if err != nil {
		return err
	}
	db.BinLog.WriteInfoToBinlogInfo()
	saveAndPrintDiskData(db.Lsm)
	return nil
}

func saveAndPrintDiskData(lsmMap *map[uint16]*lsm.LSMTree) {
	for _, lsm := range *lsmMap {
		lsm.SaveActiveToDiskOnExit()
		lsm.PrintDiskDataToFile(string(lsm.LsmPath))
	}
}

func (db *XcDB) Close() {

}

// GetVersion 获取指定键的版本号
func (db *XcDB) GetVersion(key []byte) (uint32, error) {
	//// 根据数据类型进行版本号查询
	//switch db.GetDataType(key) {
	//case XCDB_String:
	//	// 根据键获取版本号
	//	return db.GetVersionForString(key)
	//case XCDB_StringSet:
	//	// 根据键获取版本号
	//	return db.GetVersionForStringSet(key)
	//case XCDB_List:
	//	// 根据键获取版本号
	//	return db.GetVersionForList(key)
	//case XCDB_Hash:
	//	// 根据键获取版本号
	//	return db.GetVersionForHash(key)
	//case XCDB_Set:
	//	// 根据键获取版本号
	//	return db.GetVersionForSet(key)
	//default:
	//	return 0, fmt.Errorf("unsupported data type")
	//}
	return 0, nil
}

// GetVersionForString 根据键获取 String 类型数据的版本号
func (db *XcDB) GetVersionForString(key []byte) (uint32, error) {
	// 实现具体的逻辑来获取 String 类型数据的版本号
	// 例如从数据库中查询版本号等操作
	return 0, nil
}

// GetVersionForStringSet 根据键获取 StringSet 类型数据的版本号
func (db *XcDB) GetVersionForStringSet(key []byte) (uint32, error) {
	// 实现具体的逻辑来获取 StringSet 类型数据的版本号
	// 例如从数据库中查询版本号等操作
	return 0, nil
}

//
//// RemoveConnection 从连接池中移除一个连接
//func (db *XcDB) RemoveConnection(conn *driverConn) {
//	db.Pool.mu.Lock()
//	defer db.Pool.mu.Unlock()
//
//	for i, c := range db.Pool.freeConn {
//		if c == conn {
//			copy(db.Pool.freeConn[i:], db.Pool.freeConn[i+1:])
//			db.Pool.freeConn[len(db.Pool.freeConn)-1] = nil
//			db.Pool.freeConn = db.Pool.freeConn[:len(db.Pool.freeConn)-1]
//			db.Pool.numOpen--
//			break
//		}
//	}
//}
