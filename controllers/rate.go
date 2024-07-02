package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lyyubava/solidgate-software-engineering-school.git/models"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type CurrencyInfo struct {
	Currency     string  `json:"cc"`
	CurrencyCode int     `json:"r030"`
	Rate         float32 `json:"rate"`
	ExchangeDate string  `json"exchangedate"`
	Txt          string  `json:"txt"`
}

type CurrencyData []CurrencyInfo

func (c CurrencyData) Get(currencyCode string) CurrencyInfo {
	for _, currencyInfo := range c {
		if currencyInfo.Currency == currencyCode {
			return currencyInfo
		}
	}
	return CurrencyInfo{}
}

func Rate(c *gin.Context) {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	rateDb := models.Rate{}
	models.DB.Model(models.Rate{}).Where("exchange_date = ?", today).First(&rateDb)
	if rateDb.ExchangeDate.UTC() == today.UTC() {
		c.JSON(http.StatusOK, gin.H{"rate": rateDb.Rate})
		return
	}

	var currencyDataResp CurrencyData
	response, err := http.Get(os.Getenv("EXCHANGERATE_API_URL"))
	defer response.Body.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exchange rate data"})
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	err = json.Unmarshal(body, &currencyDataResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal exchange rate data"})
		return
	}

	resp := currencyDataResp.Get("USD")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return

	}
	exchangeDate := strings.Split(resp.ExchangeDate, ".")
	year, _ := strconv.Atoi(exchangeDate[2])
	month, _ := strconv.Atoi(exchangeDate[1])
	day, _ := strconv.Atoi(exchangeDate[0])
	rateExchangeDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	rate := models.Rate{Rate: resp.Rate, ExchangeDate: rateExchangeDate}
	models.DB.Create(&rate)
	c.JSON(http.StatusOK, gin.H{"rate": resp.Rate})

}
