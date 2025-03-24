package main

import (
	"github.com/rebeccastewartarcherdx/receipt-processor/controller"
	"github.com/rebeccastewartarcherdx/receipt-processor/dao"
	"github.com/rebeccastewartarcherdx/receipt-processor/handler"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/buntdb"
)

var (
	d dao.ReceiptDao
	c controller.Processor
	h handler.ReceiptHandler

	port = "localhost:6790"
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
	if err = e.Start(port); err != nil {
		e.Logger.Fatal(err)
	}
}
