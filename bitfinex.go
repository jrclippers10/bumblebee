package main

import (
    "encoding/json"
    "log"
    "strconv"
    "time"
)

var bitfinexURLs map[string]string

func init() {
  bitfinexURLs = map[string]string{
    "BTC" : "https://api.bitfinex.com/v1/book/BTCUSD",
  }
}

type BitfinexOrderUnit struct{
  Price string
  Amount string
  Timestamp string
}

type BitfinexOrderBook struct{
  Bids []BitfinexOrderUnit
  Asks []BitfinexOrderUnit
}

func (b *BitfinexOrderBook) toSof() (s Sof) {
  s.Timestamp = time.Now()
  s.Bids = make([]SofUnit, len(b.Bids))
  s.Asks = make([]SofUnit, len(b.Asks))
  for i, p := range b.Bids {
    s.Bids[i].Price, _ = strconv.ParseFloat(p.Price, 64)
    s.Bids[i].Quantity, _ = strconv.ParseFloat(p.Amount, 64)
  }
  for i, p := range b.Asks {
    s.Asks[i].Price, _ = strconv.ParseFloat(p.Price, 64)
    s.Asks[i].Quantity, _ = strconv.ParseFloat(p.Amount, 64)
  }
  return
}

func newBitfinexOrderBook(b []byte) (o BitfinexOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func bitfinex() {
  b, err := getURL(bitfinexURLs["BTC"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newBitfinexOrderBook(b)
  s := o.toSof()
  log.Println(s)
}

