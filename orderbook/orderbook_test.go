package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}

}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)
	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}
func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrdwrBook()
	sellerOrderA := NewOrder(false, 10)
	// sellerOrderB := NewOrder(false, 10) 
	ob.PlaceLimitOrder(10_000, sellerOrderA)
	// ob.PlaceLimitOrder(11_000, sellerOrderB)
	// assert(t, len(ob.Orders), 2)
	// assert(t, ob.Orders[sellerOrderA.ID], sellerOrderA)
	// assert(t, ob.Orders[sellerOrderB.ID], sellerOrderB)
	// assert(t, len(ob.asks), 2)

}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrdwrBook()
	sellOrder := NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrder)
	buyOrder := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	// assert(t, len(ob.asks), 1)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeField, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, matches[0].Bid.isFiiled(), true)
	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFiile(t *testing.T) {
	ob := NewOrdwrBook()

	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)
	buyOrderD := NewOrder(true, 1)
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)
	ob.PlaceLimitOrder(5_000, buyOrderD)

	assert(t, ob.BidTotalVolume(), 24.00)
	sellOrder := NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)
	assert(t, ob.BidTotalVolume(), 4.0)
	assert(t, len(matches), 3)
	// assert(t, len(ob.bids), 1)
	// assert(t,matches[0].Ask,sellOrder)
	// assert(t,matches[0].Bid,buyOrder)
	// assert(t,matches[0].SizeField,10.0)
	// assert(t,matches[0].Price,10_000.0)
	// assert(t,matches[0].Bid.isFiiled(),true)
	fmt.Printf("%+v", matches)
}
func TestCancelOrder(t *testing.T) {
	ob := NewOrdwrBook()

	buyOrderA := NewOrder(true, 5)
	// buyOrderB := NewOrder(true, 8)
	// buyOrderC := NewOrder(true, 10)
	// buyOrderD := NewOrder(true, 1)
	// ob.PlaceLimitOrder(5_000, buyOrderC)
	// ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)
	// ob.PlaceLimitOrder(5_000, buyOrderD)

	assert(t, ob.BidTotalVolume(), 5.0)
	ob.CancelOrder(buyOrderA)

	assert(t, ob.BidTotalVolume(), 0.0)
	_,ok := ob.Orders[buyOrderA.ID]
	assert(t, ok, false)
	// sellOrder := NewOrder(false, 20)
	// matches := ob.PlaceMarketOrder(sellOrder)
	// assert(t, ob.BidTotalVolume(), 4.0)
	// assert(t, len(matches), 3)
	// assert(t,len(ob.bids),1)
	// assert(t,matches[0].Ask,sellOrder)
	// assert(t,matches[0].Bid,buyOrder)
	// assert(t,matches[0].SizeField,10.0)
	// assert(t,matches[0].Price,10_000.0)
	// assert(t,matches[0].Bid.isFiiled(),true)
	// fmt.Printf("%+v", matches)
}
  