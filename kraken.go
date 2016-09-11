package main

import (
    "encoding/json"
    "log"
    "strconv"
    "time"
)
var KrakenURLs map[string]string

func init() {
  KrakenURLs = map[string]string{
    "BTC" : "https://api.kraken.com/0/public/Depth?pair=XXBTZUSD",
  }
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

func (b *KrakenOrderBook) toSof() (s Sof) {
  s.Timestamp = time.Now()
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

func newKrakenOrderBook(b []byte) (o KrakenOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func kraken() {
  b, err := getURL(KrakenURLs["BTC"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newKrakenOrderBook(b)
  s := o.toSof()
  log.Println(s)
}
