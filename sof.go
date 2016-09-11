package main

import (
    "time"
)

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
