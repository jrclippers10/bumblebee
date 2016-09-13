package main

import (
    "encoding/json"
    "log"
    "strconv"
    "time"
)

type Kraken struct{
  Client
}

type KrakenError struct{}

type KrakenResult struct{
  Bids [][3]interface{}
  Asks [][3]interface{}
}

type KrakenOrderBook struct{
  Error []KrakenError
  Result map[string]KrakenResult
}

func (c Kraken) newOrderBook(b []byte) (o KrakenOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func (b *KrakenOrderBook) toSof() (s Sof) {
  s.Timestamp = time.Now()
  // This line will be a problem when adding more currency pairs
  r := b.Result["XXBTZUSD"]
  s.Bids = make([]SofUnit, len(r.Bids))
  s.Asks = make([]SofUnit, len(r.Asks))
  for i, p := range r.Bids {
    s.Bids[i].Price, _ = strconv.ParseFloat(p[0].(string), 64)
    s.Bids[i].Quantity, _ = strconv.ParseFloat(p[1].(string), 64)
  }
  for i, p := range r.Asks {
    s.Asks[i].Price, _ = strconv.ParseFloat(p[0].(string), 64)
    s.Asks[i].Quantity, _ = strconv.ParseFloat(p[1].(string), 64)
  }
  return
}

func (c Kraken) run() (s Sof) {
  b, err := c.getBTCUSDOrderBook()
  if err != nil {
    log.Fatal(err)
  }
  o, err := c.newOrderBook(b)
  s = o.toSof()
  log.Println(s)
  return
}
