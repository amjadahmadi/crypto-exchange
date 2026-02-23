package main

import (
	// "context"
	// "fmt"
	"net/http"
	"strconv"
	// "time"

	// "github.com/amjadahmadi/crypto-exchange/config"
	"github.com/amjadahmadi/crypto-exchange/orderbook"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	ex := NewExchange()
	e.GET("/book/:market", ex.handleBook)
	e.POST("/order", ex.handlerPlaceOrder)
	e.DELETE("/order/:id", ex.CancelOrder)
	e.Start(":3000")

}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrdwrBook()
	return &Exchange{
		orderbooks: orderbooks,
	}
}

type PlaceOrderRequest struct {
	Type   OrderType
	Bid    bool
	Size   float64
	Price  float64
	Market Market
}
type OrderD struct {
	ID        int64
	Price     float64
	Size      float64
	Bid       bool
	TimeStamp int64
}
type OrderbookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*OrderD
	Bids           []*OrderD
}

func (ex *Exchange) handleBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]

	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "Market not found"})
	}
	orderbokkData := OrderbookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*OrderD{},
		Bids:           []*OrderD{},
	}
	for _, limit := range ob.Asks() {
		for _, orders := range limit.Orders {
			o := OrderD{
				ID:        orders.ID,
				Price:     limit.Price,
				Size:      orders.Size,
				Bid:       orders.Bid,
				TimeStamp: orders.TimeStamp,
			}
			orderbokkData.Asks = append(orderbokkData.Asks, &o)
		}
	}
	for _, limit := range ob.Bids() {
		for _, orders := range limit.Orders {
			o := OrderD{
				ID:        orders.ID,
				Price:     limit.Price,
				Size:      orders.Size,
				Bid:       orders.Bid,
				TimeStamp: orders.TimeStamp,
			}
			orderbokkData.Bids = append(orderbokkData.Bids, &o)
		}
	}
	return c.JSON(http.StatusOK, orderbokkData)
}

type CancelOrderRequest struct {
	ID  int64
	Bid bool
}

func (ex *Exchange) CancelOrder(c echo.Context) error {
	idstr := c.Param("id")
	id, _ := strconv.Atoi(idstr)
	ob := ex.orderbooks[MarketETH]
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	return c.JSON(http.StatusBadRequest, map[string]any{"msg": "canceled"})
}

type MatchedOrder struct {
	Price float64
	Size  float64
	ID    int64
}

func (ex *Exchange) handlerPlaceOrder(c echo.Context) error {
	placeOrder := PlaceOrderRequest{}
	err := c.Bind(&placeOrder)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error")
	}




	// address1 := Address{"1 Lakewood Way", "Elwood City", "PA"}
	// student1 := Student{FirstName: "Arthur", Address: address1, Age: 8}

	market := Market(placeOrder.Market)
	ob := ex.orderbooks[market]
	order := orderbook.NewOrder(placeOrder.Bid, placeOrder.Size)

	isBid := false
	if order.Bid {
		isBid = true
	}
	if placeOrder.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrder.Price, order)
		return c.JSON(200, map[string]any{"msg": "order placed"})
	}
	if placeOrder.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		matchedOrder := make([]*MatchedOrder, len(matches))
		for i := 0; i < len(matchedOrder); i++ {
			id := matches[i].Bid.ID
			if isBid {
				id = matches[i].Ask.ID
			}
			matchedOrder[i] = &MatchedOrder{
				ID:    id,
				Size:  matches[i].SizeField,
				Price: matches[i].Price,
			}
		}
		return c.JSON(200, map[string]any{"matches": len(matches)})
	}
	return nil
}
