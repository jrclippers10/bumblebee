package main

import (
    "strconv"
    "time"
)

type BitstampOrderBook struct{
  Timestamp string `json:"timestamp"`
  Bids [][2]string `json:"bids"`
  Asks [][2]string `json:"asks"`
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

type KrakenError struct{}

type KrakenResult struct{
  Bids [][3]interface{}
  Asks [][3]interface{}
}

type KrakenOrderBook struct{
  Error []KrakenError
  Result map[string]KrakenResult
}

// Standardized Orderbook Format
type SofUnit struct{
  Price float64
  Quantity float64
}

type Sof struct{
  Timestamp time.Time
  Bids []SofUnit
  Asks []SofUnit
}

// Descending
type SortBids []SofUnit
func (s SortBids) Len() int {
    return len(s)
}
func (s SortBids) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s SortBids) Less(i, j int) bool {
    return s[i].Price > s[j].Price
}

// Ascending
type SortAsks []SofUnit
func (s SortAsks) Len() int {
    return len(s)
}
func (s SortAsks) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s SortAsks) Less(i, j int) bool {
    return s[i].Price < s[j].Price
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