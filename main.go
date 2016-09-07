package main

import (
    "net/http"
    "io/ioutil"
    "time"
    // "errors"
    "log"
    // "sort"
    // "sync"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    // "net/http/httputil"
)

var bitstampMessages chan *http.Response
var krakenMessages chan *http.Response
var bitfinexMessages chan *http.Response
var bitcoinToUSD map[string]string

type bitstampOrderBook struct{
  Timestamp string `json:"timestamp"`
  Bids [][2]string `json:"bids"`
  Asks [][2]string `json:"asks"`
}

type bitfinexOrderUnit struct{
  Price string
  Amount string
  Timestamp string
}

type bitfinexOrderBook struct{
  Bids []bitfinexOrderUnit
  Asks []bitfinexOrderUnit
}

type krakenError struct{}

type krakenResult struct{
  Bids [][3]interface{}
  Asks [][3]interface{}
}

type krakenOrderBook struct{
  Error []krakenError
  Result map[string]krakenResult
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
type sortBids []SofUnit
func (s sortBids) Len() int {
    return len(s)
}
func (s sortBids) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s sortBids) Less(i, j int) bool {
    return s[i].Price > s[j].Price
}

// Ascending
type sortAsks []SofUnit
func (s sortAsks) Len() int {
    return len(s)
}
func (s sortAsks) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s sortAsks) Less(i, j int) bool {
    return s[i].Price < s[j].Price
}


func (b *bitstampOrderBook) toSof() (s Sof) {
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

func (b *bitfinexOrderBook) toSof() (s Sof) {
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

func (b *krakenOrderBook) toSof() (s Sof) {
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

func init() {
  bitstampMessages = make(chan *http.Response)
  krakenMessages = make(chan *http.Response)
  bitfinexMessages = make(chan *http.Response)

  bitcoinToUSD = map[string]string{
    "bitfinex" : "https://api.bitfinex.com/v1/book/BTCUSD",
    "bitstamp" : "https://www.bitstamp.net/api/order_book/",
    "kraken" : "https://api.kraken.com/0/public/Depth?pair=XXBTZUSD",
  }
}

func main() {
  // start collectors
  // start server
  getAllExchanges()
  startCollectors()
  startServer()

}

func newBitstampOrderBook(b []byte) (o bitstampOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func newBitfinexOrderBook(b []byte) (o bitfinexOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func newKrakenOrderBook(b []byte) (o krakenOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func getURL(s string) (body []byte, err error) {
  req, err := http.NewRequest("GET", s, nil)
  req.Header.Set("User-Agent", "Bumblebee v1.0")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      return
  }
  defer resp.Body.Close()
  body, err = ioutil.ReadAll(resp.Body)
  return
}

func startCollectors() {
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

func startServer() {
    r := mux.NewRouter()
    r.HandleFunc("/", IndexHandler)
    log.Fatal(http.ListenAndServe(":8000", r))
}


func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Bitcoin!\n"))
}

func getAllExchanges() {
    bitfinex()
    bitstamp()
    kraken()
}

func bitstamp() {
  b, err := getURL(bitcoinToUSD["bitstamp"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newBitstampOrderBook(b)
  s := o.toSof()
  log.Println(s)
}

func bitfinex() {
  b, err := getURL(bitcoinToUSD["bitfinex"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newBitfinexOrderBook(b)
  s := o.toSof()
  log.Println(s)
}

func kraken() {
  b, err := getURL(bitcoinToUSD["kraken"])
  if err != nil {
    log.Fatal(err)
  }
  o, err := newKrakenOrderBook(b)
  s := o.toSof()
  log.Println(s)
}

