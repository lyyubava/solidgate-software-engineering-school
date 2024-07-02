package main

import (
	"fmt"
	"github.com/lyyubava/solidgate-software-engineering-school.git/controllers"
	"github.com/lyyubava/solidgate-software-engineering-school.git/models"
	"github.com/robfig/cron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

var DB *gorm.DB

type Message struct {
	Rate float32
	Date time.Time
}

func Database(connString string) {
	database, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic("Fail to connect to db")
	}
	DB = database
}

func GetAllEmails() (emailSli []string) {
	rows, err := DB.Model(&models.Email{}).Rows()
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	for rows.Next() {
		var email models.Email
		DB.ScanRows(rows, &email)
		emailSli = append(emailSli, email.Email)
	}
	return
}
func SendEmail(msgBody string, mailTo string, subject string) error {
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")

	msg := fmt.Sprintf("From: %s\n To: %s\nSubject: %s\n\n %s", from, mailTo, subject, msgBody)

	err := smtp.SendMail(os.Getenv("SMTP_URI"),
		smtp.PlainAuth("", from, password, os.Getenv("SMTP_HOST")),
		from, []string{mailTo}, []byte(msg))

	if err != nil {
		log.Printf("error while sending email: %s", err)
		return err

	}

	log.Printf("email was successfully sent to %s", mailTo)
	return nil
}

func SendEmails(msgBody string, emails []string, subject string) {
	for _, email := range emails {
		SendEmail(msgBody, email, subject)
	}
}

func getRateToday() (models.Rate, error) {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	rateDb := models.Rate{}
	DB.Model(models.Rate{}).Where("exchange_date = ?", today).First(&rateDb)
	if rateDb.ExchangeDate.UTC() == today.UTC() {
		return rateDb, nil
	}

	var currencyDataResp controllers.CurrencyData
	resp := currencyDataResp.Get("USD")
	response, err := http.Get(os.Getenv("EXCHANGERATE_API_URL"))
	if err != nil {
		return rateDb, err
	}
	defer response.Body.Close()
	exchangeDate := strings.Split(resp.ExchangeDate, ".")
	year, _ := strconv.Atoi(exchangeDate[2])
	month, _ := strconv.Atoi(exchangeDate[1])
	day, _ := strconv.Atoi(exchangeDate[0])
	rateExchangeDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	rate := models.Rate{Rate: resp.Rate, ExchangeDate: rateExchangeDate}
	DB.Create(&rate)
	fmt.Println(rateExchangeDate)
	return rate, nil

}

func (msg Message) formMessageBody() (msgBody string) {
	msgBody = "US Dollar Exchange Rate for " + msg.Date.Format("2006-01-02") + " is " + strconv.FormatFloat(float64(msg.Rate), 'f', -1, 64)
	return
}

func (msg Message) formMessageSubject() (msgSubject string) {
	msgSubject = "US Dollar Exchange Rate for " + msg.Date.Format("2006-01-02")
	return
}

func main() {

	connString := os.Getenv("DATABASE_CONNECTION_STRING")
	Database(connString)
	fmt.Println("connString: ", connString)
	emails := GetAllEmails()
	c := cron.New()
	err := c.AddFunc("* * * * *", func() {
		rate, err := getRateToday()
		fmt.Println(rate.Rate, rate.ExchangeDate, err)
		if err == nil {
			message := Message{Rate: rate.Rate, Date: rate.ExchangeDate}

			SendEmails(message.formMessageBody(), emails, message.formMessageSubject())
		} else {
			panic(err)
		}

	})
	c.Run()

	if err != nil {
		panic(err)
	}

}
