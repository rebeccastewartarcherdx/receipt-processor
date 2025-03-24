package controller

import (
	"testing"
	"time"

	"github.com/rebeccastewartarcherdx/receipt-processor/models"
)

func Test_convertToCents(t *testing.T) {
	type args struct {
		totalString string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{totalString: "31.23"},
			want:    3123,
			wantErr: false,
		},
		{
			name:    "invalid number",
			args:    args{totalString: "invalid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToCents(tt.args.totalString)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToCents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertToCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateDescriptionPoints(t *testing.T) {
	type args struct {
		items []models.Item
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "divisble by 3",
			args: args{items: []models.Item{
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            1225,
				},
			}},
			want: 3,
		},
		{
			name: "not divisble by 3",
			args: args{items: []models.Item{
				{
					ShortDescription: "Emils Cheese Pizza!",
					Price:            1225,
				},
			}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDescriptionPoints(tt.args.items); got != tt.want {
				t.Errorf("calculateDescriptionPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_alphaNumericCount(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "all alpha",
			args: args{input: "abcd"},
			want: 4,
		},
		{
			name: "all numeric",
			args: args{input: "123"},
			want: 3,
		},
		{
			name: "alphanumeric",
			args: args{input: "abc123!@#"},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alphaNumericCount(tt.args.input); got != tt.want {
				t.Errorf("alphaNumericCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculatePoints(t *testing.T) {
	layout := "2006-01-02 15:04"
	outWindowEvenDate, _ := time.Parse(layout, "2022-01-02 13:13")
	outWindowOddDate, _ := time.Parse(layout, "2022-01-01 13:13")
	inWindowEvenDate, _ := time.Parse(layout, "2022-01-02 14:13")
	type args struct {
		receipt *models.Receipt
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "simple",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowEvenDate,
				Total:            126,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            126,
					},
				},
			}},
			want: 6,
		},
		{
			name: "no cents",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowEvenDate,
				Total:            100,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            100,
					},
				},
			}},
			want: 81,
		},
		{
			name: "multiple of 25",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowEvenDate,
				Total:            125,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            125,
					},
				},
			}},
			want: 31,
		},
		{
			name: "2 items in receipt",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowEvenDate,
				Total:            252,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            126,
					},
					{
						ShortDescription: "Diet Coke - 12-oz",
						Price:            126,
					},
				},
			}},
			want: 11,
		},
		{
			name: "multiple of 3 description",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowEvenDate,
				Total:            126,
				Items: []models.Item{
					{
						ShortDescription: "Coke - 12-oz",
						Price:            126,
					},
				},
			}},
			want: 7,
		},
		{
			name: "odd day",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: outWindowOddDate,
				Total:            126,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            126,
					},
				},
			}},
			want: 12,
		},
		{
			name: "time between 2 and 4",
			args: args{receipt: &models.Receipt{
				Retailer:         "Target",
				PurchaseDateTime: inWindowEvenDate,
				Total:            126,
				Items: []models.Item{
					{
						ShortDescription: "Pepsi - 12-oz",
						Price:            126,
					},
				},
			}},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePoints(tt.args.receipt); got != tt.want {
				t.Errorf("calculatePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
