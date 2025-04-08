package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"superfrete-api/model"
	"superfrete-api/repository"

	"github.com/gin-gonic/gin"
)

type QuoteHandler struct {
	repository repository.QuoteRepository
}

func NewQuoteHandler(repo repository.QuoteRepository) QuoteHandler {
	return QuoteHandler{
		repository: repo,
	}
}

func (q *QuoteHandler) GetLastQuote(ctx *gin.Context) {
	/* quotes := []model.Quote{ //TODO: Deixar para testes
		{
			ID:       1,
			Name:     "EXPRESSO FR",
			Service:  "Rodoviário",
			Deadline: 3,
			Price:    17,
		},
		{
			ID:       2,
			Name:     "Correios",
			Service:  "SEDEX",
			Deadline: 1,
			Price:    20.99,
		},
	} */

	limit, err := strconv.Atoi(ctx.DefaultQuery("last_quotes", "-1"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "limite inválido"})
		return
	}

	lastQuotes, err := q.repository.GetLastQuotes(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var response model.CarrierResponse
	response.CalculateQuoteMetrics(lastQuotes)

	ctx.JSON(http.StatusOK, response)
}

func (q *QuoteHandler) PostQuote(ctx *gin.Context) {
	quotes, err := q.callFRAPI(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = q.saveQuote(quotes)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, quotes)
}

func (q *QuoteHandler) saveQuote(quote []model.Quote) error {
	err := q.repository.CreateQuote(quote)
	if err != nil {
		return err
	}

	return nil
}

func (q *QuoteHandler) callFRAPI(ctx *gin.Context) ([]model.Quote, error) {
	var quoteRequest model.QuoteRequest

	err := ctx.ShouldBindJSON(&quoteRequest)
	if err != nil {
		return nil, err
	}

	cnpj := os.Getenv("CNPJ")
	token := os.Getenv("FRETE_RAPIDO_TOKEN")

	shippingRequest := model.ConvertQuoteToShipping(quoteRequest)
	shippingRequest.Shipper.RegisteredNumber = cnpj
	shippingRequest.Shipper.Token = token
	shippingRequest.Shipper.PlatformCode = os.Getenv("COD_PLATAFORMA")

	jsonData, err := json.Marshal(shippingRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://sp.freterapido.com/api/v3/quote/simulate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	freteResponse := model.FreteRapidoResponse{}
	err = json.Unmarshal(res, &freteResponse)
	if err != nil {
		return nil, err
	}

	quotes, err := freteResponse.ConvertToQuote()
	if err != nil {
		return nil, err
	}

	return quotes, nil
}
