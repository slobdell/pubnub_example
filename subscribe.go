package main

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
)

func startSubscription(pubnubChannelId string, successChannel, errorChannel chan []byte) {
	pubnub := messaging.NewPubnub(
		PUBLISH_KEY,
		SUBSCRIBE_KEY,
		SECRET_KEY,
		CIPHER_KEY,
		USE_SSL,
		"",  // custom UUID
		nil, // optional logger
	)
	go pubnub.Subscribe(
		pubnubChannelId,
		"", // no idea
		successChannel,
		false,
		errorChannel,
	)
}

func infinitePubnubRead(rawMessages chan string, successChannel, errorChannel chan []byte) {
	var incomingMessage []interface{}
	for {
		select {
		case response := <-successChannel:

			err := json.Unmarshal(response, &incomingMessage)
			if err != nil {
				fmt.Println(err)
				continue
			}

			switch messageData := incomingMessage[0].(type) {
			case float64:
				// this case reached on startup, not exactly sure what this is all about
			case []interface{}:
				rawMessages <- messageData[0].(string)
			default:
				panic(fmt.Sprintf("Unknown type: %T", messageData))
			}

		case err := <-errorChannel:
			fmt.Println(string(err))
		case <-messaging.SubscribeTimeout():
			fmt.Println("Subscribe() timeout")
		}
	}
}

func extractMessagesFromChannel(pubnubChannelId string, rawMessages chan string) {
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	startSubscription(pubnubChannelId, successChannel, errorChannel)
	infinitePubnubRead(rawMessages, successChannel, errorChannel)

}
