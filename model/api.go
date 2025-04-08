package model

type ShippingRequest struct {
	Shipper        Shipper         `json:"shipper"`
	Recipient      RecipientDetail `json:"recipient"`
	Dispatchers    []Dispatcher    `json:"dispatchers"`
	Channel        string          `json:"channel"`
	Filter         int             `json:"filter"`
	Limit          int             `json:"limit"`
	Identification string          `json:"identification"`
	Reverse        bool            `json:"reverse"`
	SimulationType []int           `json:"simulation_type"`
	Returns        Returns         `json:"returns"`
}

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type RecipientDetail struct {
	Type             int    `json:"type"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Country          string `json:"country"`
	Zipcode          int    `json:"zipcode"`
}

type Dispatcher struct {
	RegisteredNumber string   `json:"registered_number"`
	Zipcode          int      `json:"zipcode"`
	TotalPrice       float64  `json:"total_price"`
	Volumes          []Volume `json:"volumes"`
}

type Volume struct {
	Amount        int     `json:"amount"`
	AmountVolumes int     `json:"amount_volumes"`
	Category      string  `json:"category"`
	SKU           string  `json:"sku"`
	Tag           string  `json:"tag"`
	Description   string  `json:"description"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Consolidate   bool    `json:"consolidate"`
	Overlaid      bool    `json:"overlaid"`
	Rotate        bool    `json:"rotate"`
}

type Returns struct {
	Composition  bool `json:"composition"`
	Volumes      bool `json:"volumes"`
	AppliedRules bool `json:"applied_rules"`
}

type FreteRapidoResponse struct {
	Dispatchers []struct {
		Offers []struct {
			Carrier struct {
				Name string `json:"name"`
			} `json:"carrier"`
			Service      string       `json:"service"`
			FinalPrice   float64      `json:"final_price"`
			DeliveryTime deliveryTime `json:"delivery_time"`
		} `json:"offers"`
	} `json:"dispatchers"`
}

type deliveryTime struct {
	Hours int `json:"hours"`
	Days  int `json:"days"`
}

func (f *FreteRapidoResponse) ConvertToQuote() ([]Quote, error) {
	var quotes []Quote
	for _, dispatcher := range f.Dispatchers {
		for _, offer := range dispatcher.Offers {
			quote := Quote{
				Name:    offer.Carrier.Name,
				Service: offer.Service,
				Price:   offer.FinalPrice,
			}

			quote.fillDdl(offer.DeliveryTime)
			quotes = append(quotes, quote)
		}
	}

	return quotes, nil
}
