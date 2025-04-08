package model

type CarrierResponse struct {
	Carriers []QuoteMetrics `json:"carriers"`
}

type QuoteMetrics struct {
	Metrics       []CarrierMetrics `json:"metrics"`
	CheapestQuote Quote            `json:"cheapest_quote"`
	HigherQuote   Quote            `json:"higher_quote"`
}

type CarrierMetrics struct {
	Name         string  `json:"name"`
	Count        int     `json:"count"`
	TotalPrice   float64 `json:"total_price"`
	AveragePrice float64 `json:"average_price"`
}

func (res *CarrierResponse) CalculateQuoteMetrics(quotes []Quote) {
	carrierMap := make(map[string]*CarrierMetrics)
	var cheapest, mostExpensive Quote
	first := true

	for _, quote := range quotes {
		if first || quote.Price < cheapest.Price {
			cheapest = quote
		}
		if first || quote.Price > mostExpensive.Price {
			mostExpensive = quote
		}
		first = false

		if _, exists := carrierMap[quote.Name]; !exists {
			carrierMap[quote.Name] = &CarrierMetrics{
				Name:       quote.Name,
				Count:      1,
				TotalPrice: quote.Price,
			}
		} else {
			carrierMap[quote.Name].Count++
			carrierMap[quote.Name].TotalPrice += quote.Price
		}
	}

	var carrierMetrics []CarrierMetrics
	for _, m := range carrierMap {
		m.AveragePrice = m.TotalPrice / float64(m.Count)
		carrierMetrics = append(carrierMetrics, *m)
	}

	res.Carriers = []QuoteMetrics{
		{
			Metrics:       carrierMetrics,
			CheapestQuote: cheapest,
			HigherQuote:   mostExpensive,
		},
	}
}
