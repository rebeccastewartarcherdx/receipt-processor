package dao

import (
	"receipt_processor/models"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
)

type (
	ReceiptDao interface {
		Insert(receipt *models.Receipt) (string, error)
		Get(id string) (receipt *models.Receipt, err error)
	}

	receiptDao struct {
		db *buntdb.DB
	}
)

func NewReceiptDao(db *buntdb.DB) ReceiptDao {
	return &receiptDao{
		db: db,
	}
}

func (d *receiptDao) Insert(receipt *models.Receipt) (string, error) {
	id := uuid.NewString()
	b, err := json.Marshal(receipt)
	if err != nil {
		return "", err
	}
	err = d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(id, string(b), nil)
		return err
	})
	return id, err
}

func (d *receiptDao) Get(id string) (receipt *models.Receipt, err error) {
	err = d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		receipt = &models.Receipt{}
		err = json.Unmarshal([]byte(val), receipt)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}
