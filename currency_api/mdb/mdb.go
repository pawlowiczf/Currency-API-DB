package mdb

import (
	"currency_api/api"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/mattn/go-sqlite3"
)

type DailyExchangeDetailsEntry struct {
	FromCurrency  string
	ToCurrency    string
	LastRefreshed string
	ExchangeRate  string
}

func CreateDatabase(db *sql.DB) {
	//
	_, err := db.Exec(`
		CREATE TABLE dailyExchange (
			id INTEGER PRIMARY KEY, 
			fromCurrency TEXT, 
			toCurrency TEXT,
			lastRefreshed TEXT, 
			exchangeRate TEXT
		)
	`)

	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX idx_from_to_date ON dailyExchange(fromCurrency, toCurrency, lastRefreshed)
	`)

	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func exchangeDetailsEntriesFromRow(rows *sql.Rows) ([]DailyExchangeDetailsEntry, error) {
	//
	exchangeEntries := []DailyExchangeDetailsEntry{}

	var fromCurrency string
	var toCurrency string
	var lastRefreshed string
	var exchangeRate string

	defer rows.Close()

	for rows.Next() {

		status := rows.Scan(&fromCurrency, &toCurrency, &lastRefreshed, &exchangeRate)

		if status != nil {
			log.Println(status)
			return nil, status
		}

		exchangeEntry := DailyExchangeDetailsEntry{FromCurrency: fromCurrency, ToCurrency: toCurrency, LastRefreshed: lastRefreshed, ExchangeRate: exchangeRate}
		exchangeEntries = append(exchangeEntries, exchangeEntry)
	}

	return exchangeEntries, nil
}

func CreateDailyExchangeDetailsEntry(db *sql.DB, exchangeDetails api.DailyExchangeDetails, fromCurrency, toCurrency string) error {
	//
	parsedTime, err := time.Parse("2006-01-02", exchangeDetails.Date)
	newParsedTime := parsedTime.String()[:10]

	if err != nil {
		log.Fatalln(err)
		return err
	}

	if exchangeDetails.Currency[toCurrency] <= 0 {
		return nil
	}
	
	target := DailyExchangeDetailsEntry{}
	err = db.QueryRow(`SELECT fromCurrency, toCurrency, lastRefreshed, exchangeRate FROM dailyExchange WHERE lastRefreshed=? AND fromCurrency=? AND toCurrency=?`, newParsedTime, fromCurrency, toCurrency).Scan(&target.FromCurrency, &target.ToCurrency, &target.LastRefreshed, &target.ExchangeRate)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalln(err)
		}
		target.LastRefreshed = "no-date"
	}

	if target.LastRefreshed == newParsedTime && target.FromCurrency == fromCurrency && target.ToCurrency == toCurrency {
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO dailyExchange(fromCurrency, toCurrency, lastRefreshed, exchangeRate)
		VALUES(?, ?, ?, ?)
	`, fromCurrency, toCurrency, newParsedTime, strconv.FormatFloat(exchangeDetails.Currency[toCurrency], 'f', -1, 64))

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func GetDailyExchangeRateEntries(db *sql.DB, fromCurrency, toCurrency string) ([]DailyExchangeDetailsEntry, error) {
	//
	rows, err := db.Query(`
		SELECT fromCurrency, toCurrency, lastRefreshed, exchangeRate 
		FROM dailyExchange 
		WHERE fromCurrency=? AND toCurrency=?
	`, fromCurrency, toCurrency)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return exchangeDetailsEntriesFromRow(rows)
}

func GetDailyExchange(db *sql.DB) ([]DailyExchangeDetailsEntry, error) {
	//
	rows, err := db.Query(`
		SELECT fromCurrency, toCurrency, lastRefreshed, exchangeRate
		FROM DailyExchange
	`)

	if err != nil {
		log.Fatalln(err)
	}

	return exchangeDetailsEntriesFromRow(rows)
}
