package main

import (
    "time"
)

var bitfinex Bitfinex
var bitstamp Bitstamp
var kraken Kraken

func init() {
  bitfinex = Bitfinex{
    Client: Client{
      BTC: map[string]string{
        "USD" : "https://api.bitfinex.com/v1/book/BTCUSD",
      },
    },
  }
  kraken = Kraken{
    Client: Client{
      BTC: map[string]string{
        "USD" : "https://api.kraken.com/0/public/Depth?pair=XXBTZUSD",
      },
    },
  }
  bitstamp = Bitstamp{
    Client: Client{
      BTC: map[string]string{
        "USD" : "https://www.bitstamp.net/api/order_book/",
      },
    },
  }
}

func getAllExchanges() {
    bitfinex.run()
    bitstamp.run()
    kraken.run()
}

func StartCollectors() {
  getAllExchanges()
  ticker := time.NewTicker(30 * time.Second)
  quit := make(chan struct{})
  go func() {
      for {
         select {
          case <- ticker.C:
              getAllExchanges()
          case <- quit:
              ticker.Stop()
              return
          }
      }
   }()
}