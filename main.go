package main

import (
	"fmt"
	ws "github.com/gorilla/websocket"
	"github.com/jamespearly/loggly"
	"github.com/preichenberger/go-coinbasepro/v2"
	"strconv"
)

func main() {

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		var tag string
		tag = "My-Go-Demo"
		loggly_client := loggly.New(tag)

		if message.Type == "ticker" {
			err = loggly_client.EchoSend("info", "TYPE: " + message.Type + ", PRODUCT ID: " + message.ProductID + ", TRADE ID: " + strconv.Itoa(message.TradeID) + ", SEQUENCE: " + strconv.Itoa(int(message.Sequence)) + ", PRICE: " + message.Price + ", SIDE: " + message.Side + ", LAST SIZE: " + message.LastSize + ", BEST BID: " + message.BestBid + ", BEST ASK: "+ message.BestAsk)
			fmt.Println("err:", err)
		} else {
			err = loggly_client.EchoSend("error", "MESSAGE TYPE OTHER THAN TICKER RECEIVED")
		}
	}
}
