package main

import (
	"currency_api/api"
	"currency_api/mdb"
	"database/sql"
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
}

func main() {
	//
	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "cantor.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}

	db, err := sql.Open("sqlite3", args.DbPath)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	mdb.CreateDatabase(db)
	
	handleRequest(db)
}

func handleRequest(db *sql.DB) {
	//
	fromCurrency, toCurrency := "btc", "usd"
	dailyExchangeDetails, err := api.GetDailyExchangeRate(fromCurrency)

	err = mdb.CreateDailyExchangeDetailsEntry(db, dailyExchangeDetails, fromCurrency, toCurrency)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := mdb.GetDailyExchangeRateEntries(db, fromCurrency, toCurrency)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(entries)
	fmt.Println( mdb.GetDailyExchange(db) )
}