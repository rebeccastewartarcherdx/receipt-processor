package handler

import (
	"fmt"
	"receipt_processor/controller"
	"receipt_processor/models"

	"github.com/labstack/echo/v4"
)

type ReceiptHandler struct {
	controller controller.Processor
}

func NewHandler(controller controller.Processor) ReceiptHandler {
	return ReceiptHandler{controller: controller}
}

func (h ReceiptHandler) ProcessReceipt(c echo.Context) error {
	receipt := models.UnprocessedReceipt{}
	if err := c.Bind(&receipt); err != nil {
		return c.JSON(400, "invalid request")
	}

	id, err := h.controller.Process(c.Request().Context(), &receipt)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, struct {
		Id string `json:"id"`
	}{Id: id})
}

func (h ReceiptHandler) GetPoints(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(400, "id required")
	}
	points, err := h.controller.RetrievePoints(c.Request().Context(), id)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, struct {
		Points string `json:"points"`
	}{Points: fmt.Sprintf("%d", points)})
}
