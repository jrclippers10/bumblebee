package main

import (
    "time"
)

func StartCollectors() {
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