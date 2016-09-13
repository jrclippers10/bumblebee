package main

import (
    "sync"
    "time"
)

var (
  bitfinex Bitfinex
  bitstamp Bitstamp
  kraken Kraken
  sofChan chan Sof
)

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
  sofChan = make(chan Sof)
}

func getAllExchanges() {
    var wg sync.WaitGroup
    wg.Add(3)
    go func() {
        defer wg.Done()
        sofChan <- bitfinex.run()
    }()
    go func() {
        defer wg.Done()
        sofChan <- bitstamp.run()
    }()
    go func() {
        defer wg.Done()
        sofChan <- kraken.run()
    }()
    wg.Wait()
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