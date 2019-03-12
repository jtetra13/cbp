package main

import (
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	ws "github.com/gorilla/websocket"
	"github.com/jamespearly/loggly"
	"github.com/preichenberger/go-coinbasepro/v2"
	"os"
	"strconv"
)

type Ticker struct	{
	Type string `json:"type"`
	ProductID string `json:"product_id"`
	TradeID int `json:"trade_id,number"`
	Time coinbasepro.Time `json:"time,string"`
	Sequence int64 `json:"sequence,number"`
	Price string `json:"price"`
	Side string `json:"side"`
	LastSize string `json:"last_size"`
	BestBid string `json:"best_bid"`
	BestAsk string `json:"best_ask"`
}

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	var wsDialer = ws.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
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

		item := Ticker{
			Type: message.Type,
			ProductID: message.ProductID,
			TradeID: message.TradeID,
			Time: message.Time,
			Sequence: message.Sequence,
			Price: message.Price,
			Side: message.Side,
			LastSize: message.LastSize,
			BestBid: message.BestBid,
			BestAsk: message.BestAsk,
		}

		av, err := dynamodbattribute.MarshalMap(item)

		input := &dynamodb.PutItemInput{
			Item: av,
			TableName: aws.String("jbit"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println(message.Time)
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Successfully added!")

		var tag string
		tag = "My-Go-Demo"
		logglyClient := loggly.New(tag)

		if message.Type == "ticker" {
			err = logglyClient.EchoSend("info", "TYPE: " + message.Type + ", PRODUCT ID: " + message.ProductID + ", TRADE ID: " + strconv.Itoa(message.TradeID) + ", SEQUENCE: " + strconv.Itoa(int(message.Sequence)) + ", PRICE: " + message.Price + ", SIDE: " + message.Side + ", LAST SIZE: " + message.LastSize + ", BEST BID: " + message.BestBid + ", BEST ASK: "+ message.BestAsk)
			//fmt.Println("err:", err)
		} else {
			err = logglyClient.EchoSend("error", "MESSAGE TYPE OTHER THAN TICKER RECEIVED")
		}
	}

}

