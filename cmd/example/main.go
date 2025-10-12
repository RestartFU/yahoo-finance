package main

import (
	"fmt"
	"log"

	"github.com/restartfu/yahoo-finance/yahoofinance"
)

func main() {
	subscriber := yahoofinance.NewSubscriber(yahoofinance.SymbolXEQT)

	err := subscriber.Subscribe()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		poll, err := subscriber.Poll(yahoofinance.SymbolXEQT)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", poll)
	}
}
