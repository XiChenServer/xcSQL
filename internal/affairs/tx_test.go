package affairs

//import (
//	"SQL/internal/database"
//	"SQL/internal/pool"
//	"testing"
//	"time"
//)
//
//func Test_tx(t *testing.T) {
//	pool := pool.NewConnectionPool(2, 10, 30*time.Minute)pool
//	db, err := database.NewXcDB("123",pool)
//	if err != nil {
//		t.Errorf(err.Error())
//		return
//	}
//	tx := db.TX.BeginTransaction()
//
//	tx.TxGet([]byte("123"))
//
//}
