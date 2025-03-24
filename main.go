package main

import (
	"receipt_processor/controller"
	"receipt_processor/dao"
	"receipt_processor/handler"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/buntdb"
)

var (
	d dao.ReceiptDao
	c controller.Processor
	h handler.ReceiptHandler
)

func main() {
	e := echo.New()

	inMemDb, err := buntdb.Open(":memory:")
	if err != nil {
		e.Logger.Fatal("unable to init in memory db")
	}
	d = dao.NewReceiptDao(inMemDb)
	c = controller.NewProcessor(d)
	h = handler.NewHandler(c)

	e.POST("/receipts/process", h.ProcessReceipt)
	e.GET("/receipts/:id/points", h.GetPoints)
	if err = e.Start("localhost:6790"); err != nil {
		e.Logger.Fatal(err)
	}
}
