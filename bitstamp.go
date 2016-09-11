package main

import (
    "encoding/json"
    "log"
    "strconv"
    "time"
)

var BitstampURLs map[string]string

func init() {
  BitstampURLs = map[string]string{
    "BTC" : "https://www.bitstamp.net/api/order_book/",
  }
}

type BitstampOrderBook struct{
  Timestamp string `json:"timestamp"`
  Bids [][2]string `json:"bids"`
  Asks [][2]string `json:"asks"`
}

func (b *BitstampOrderBook) toSof() (s Sof) {
  s.Timestamp, _ = time.Parse(time.UnixDate, b.Timestamp)
  s.Bids = make([]SofUnit, len(b.Bids))
  s.Asks = make([]SofUnit, len(b.Asks))
  for i, p := range b.Bids {
    s.Bids[i].Price, _ = strconv.ParseFloat(p[0], 64)
    s.Bids[i].Quantity, _ = strconv.ParseFloat(p[1], 64)
  }
  for i, p := range b.Asks {
    s.Asks[i].Price, _ = strconv.ParseFloat(p[0], 64)
    s.Asks[i].Quantity, _ = strconv.ParseFloat(p[1], 64)
  }
  return
}

func newBitstampOrderBook(b []byte) (o BitstampOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func bitstamp() {
  b, err := getURL(BitstampURLs["BTC"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newBitstampOrderBook(b)
  s := o.toSof()
  log.Println(s)
}
