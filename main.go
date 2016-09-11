package main

import (
    "net/http"
    "io/ioutil"
    // "errors"
    // "log"
    // "sort"
    // "sync"
    // "encoding/json"
)

func init() {

}

func main() {
  // start collectors
  // start server
  getAllExchanges()
  StartCollectors()
  StartServer()

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

