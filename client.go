package main

import (
    "io/ioutil"
    "net/http"
)

type Client struct{
  BTC map[string]string
}

func (c Client) getBTCUSDOrderBook() (body []byte, err error) {
  return c.get(c.BTC["USD"])
}

func (c Client) get(s string) (body []byte, err error) {
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
