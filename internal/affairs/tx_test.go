package affairs

import (
	"SQL/internal/database"
	"testing"
)

func Test_tx(t *testing.T) {

	db, err := database.NewXcDB("123")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	tx := db.TX.BeginTransaction()

	tx.TxGet([]byte("123"))

}
