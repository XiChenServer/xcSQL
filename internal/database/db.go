package database

import (
	"SQL/internal/log"
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/internal/wal"
	"SQL/logs"
	"fmt"
	"sync"
)

type XcDB struct {
	StorageManager *storage.StorageManager
	Lsm            *map[uint16]*lsm.LSMTree //lsm书
	Wal            *wal.WAL                 // redo.log之类的，有关事务的
	BinLog         *log.BinlogFile
	// 读写锁，用于并发读写控制
	Mu sync.RWMutex
}

func NewXcDB(name string) (*XcDB, error) {
	var lsmMap = make(map[uint16]*lsm.LSMTree)
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
	lsmString := lsm.NewLSMTree(16, 10000, model.XCDB_String, name)
	lsmList := lsm.NewLSMTree(16, 10000, model.XCDB_List, name)
	lsmHash := lsm.NewLSMTree(16, 10000, model.XCDB_Hash, name)
	lsmSet := lsm.NewLSMTree(16, 10000, model.XCDB_Set, name)
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

	binlog, err := log.NewBinlogFile(name, 0.5*1024*1024, 5)

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
	fmt.Println("dfsd")
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
