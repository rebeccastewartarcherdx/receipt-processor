package controller

import (
	"context"
	"math"
	"receipt_processor/dao"
	"receipt_processor/models"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type (
	Processor interface {
		Process(ctx context.Context, unprocessed *models.UnprocessedReceipt) (string, error)
		RetrievePoints(ctx context.Context, id string) (int, error)
	}

	processor struct {
		dao dao.ReceiptDao
	}
)

func NewProcessor(dao dao.ReceiptDao) Processor {
	return &processor{dao: dao}
}

func (p *processor) Process(ctx context.Context, unprocessed *models.UnprocessedReceipt) (string, error) {
	// convert model to struct containing more useful data types
	receipt, err := convertModel(unprocessed)
	if err != nil {
		return "", err
	}

	// get points
	points := calculatePoints(receipt)
	receipt.Points = points

	// insert
	id, err := p.dao.Insert(receipt)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (p *processor) RetrievePoints(ctx context.Context, id string) (int, error) {
	receipt, err := p.dao.Get(id)
	if err != nil {
		return 0, err
	}
	return receipt.Points, nil
}

func convertModel(unprocessed *models.UnprocessedReceipt) (*models.Receipt, error) {
	receipt := models.Receipt{
		Retailer: unprocessed.Retailer,
	}

	// convert total to integer (cents)
	intTotal, err := convertToCents(unprocessed.Total)
	if err != nil {
		return nil, err
	}
	receipt.Total = intTotal

	// convert date time
	layout := "2006-01-02 15:04"
	date, err := time.Parse(layout, unprocessed.PurchaseDate+" "+unprocessed.PurchaseTime)
	if err != nil {
		return nil, err
	}
	receipt.PurchaseDateTime = date

	// convert items
	items := make([]models.Item, len(unprocessed.Items))
	for i, item := range unprocessed.Items {
		intPrice, err := convertToCents(item.Price)
		if err != nil {
			return nil, err
		}
		items[i] = models.Item{
			ShortDescription: item.ShortDescription,
			Price:            intPrice,
		}
	}
	receipt.Items = items
	return &receipt, nil
}

func calculatePoints(receipt *models.Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	points += alphaNumericCount(receipt.Retailer)

	// 50 points if the total is a round dollar amount with no cents.
	if receipt.Total%100 == 0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	// ASSUMPTION: this will add for round amounts in addition to the above. could put this in an else if that is not the desired behavior
	// (e.g. 20.00 would get both 50 and 25 added to the total)
	if receipt.Total%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += len(receipt.Items) / 2 * 5 // integer division

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	points += calculateDescriptionPoints(receipt.Items)

	// 6 points if the day in the purchase date is odd.
	if receipt.PurchaseDateTime.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour := receipt.PurchaseDateTime.Hour()
	if hour >= 14 && hour <= 16 {
		points += 10
	}

	return points
}

func alphaNumericCount(input string) int {
	count := 0
	for _, rune := range input {
		if unicode.IsLetter(rune) || unicode.IsNumber(rune) {
			count += 1
		}
	}
	return count
}

func calculateDescriptionPoints(items []models.Item) int {
	points := 0
	for _, item := range items {
		str := strings.TrimSpace(item.ShortDescription)
		if utf8.RuneCountInString(str)%3 == 0 {
			points += int(math.Ceil(float64(item.Price) / 100 * 0.2))
		}
	}
	return points
}

func convertToCents(totalString string) (int, error) {
	total, err := strconv.ParseFloat(totalString, 64)
	if err != nil {
		return 0, err
	}
	return int(total * 100), nil
}
