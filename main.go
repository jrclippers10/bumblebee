package main

import (
    "net/http"
    "io/ioutil"
    // "errors"
    "log"
    // "sort"
    // "sync"
    "encoding/json"
)

var bitcoinToUSD map[string]string

func init() {
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
  StartCollectors()
  StartServer()

}

func newBitstampOrderBook(b []byte) (o BitstampOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func newBitfinexOrderBook(b []byte) (o BitfinexOrderBook, err error) {
  err = json.Unmarshal(b, &o)
  if err != nil {
    log.Println("Error Unmarshaling to JSON", err)
  }
  return o, err
}

func newKrakenOrderBook(b []byte) (o KrakenOrderBook, err error) {
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

