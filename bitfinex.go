package main

import (
    "encoding/json"
    "log"
    "strconv"
    "time"
)

type Bitfinex struct{
  Client
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

func (c Bitfinex) newOrderBook(b []byte) (o BitfinexOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
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

func (c Bitfinex) run() {
  b, err := c.getBTCUSDOrderBook()
  if err != nil {
    log.Fatal(err)
  }
  o, err := c.newOrderBook(b)
  s := o.toSof()
  log.Println(s)
}

