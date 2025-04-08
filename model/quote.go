package model

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Quote struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}

type QuoteRequest struct {
	Recipient Recipient `json:"recipient"`
	Volumes   []Volumes `json:"volumes"`
}

type Recipient struct {
	Address Address `json:"address"`
}

type Address struct {
	Zipcode string `json:"zipcode"`
}

type Volumes struct {
	Category      int     `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Price         float64 `json:"price"`
	SKU           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

func ConvertQuoteToShipping(qr QuoteRequest) ShippingRequest {
	cnpj := os.Getenv("CNPJ")
	cepStr := os.Getenv("DISP_CEP")
	cep, err := strconv.Atoi(cepStr)
	if err != nil {
		log.Fatalf("Erro ao converter DISP_CEP: %v", err)
	}

	dispatcher := Dispatcher{
		RegisteredNumber: cnpj,
		Zipcode:          cep,
		TotalPrice:       0,
	}

	var totalPrice float64
	for _, v := range qr.Volumes {
		dispatcher.Volumes = append(dispatcher.Volumes, Volume{
			Amount:        v.Amount,
			AmountVolumes: v.Amount,
			Category:      fmt.Sprintf("%d", v.Category),
			SKU:           v.SKU,
			Tag:           "",
			Description:   "",
			Height:        v.Height,
			Width:         v.Width,
			Length:        v.Length,
			UnitaryPrice:  v.Price,
			UnitaryWeight: v.UnitaryWeight,
			Consolidate:   false,
			Overlaid:      false,
			Rotate:        false,
		})
		totalPrice += v.Price * float64(v.Amount)
	}
	dispatcher.TotalPrice = totalPrice

	return ShippingRequest{
		Recipient: RecipientDetail{
			Zipcode: parseZip(qr.Recipient.Address.Zipcode),
		},
		Dispatchers:    []Dispatcher{dispatcher},
		Channel:        "",
		Filter:         0,
		Limit:          0,
		Identification: "",
		Reverse:        false,
		SimulationType: []int{0},
		Returns: Returns{
			Composition:  false,
			Volumes:      false,
			AppliedRules: false,
		},
	}
}

func parseZip(zip string) int {
	i, err := strconv.Atoi(zip)
	if err != nil {
		return 0
	}
	return i
}

func (q *Quote) fillDdl(dt deliveryTime) {
	if dt.Hours != 0 {
		q.Deadline = dt.Hours
		return
	}
	q.Deadline = dt.Days
}
