package main

import (
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"os"
	"time"
)

var PUBLISH_KEY = os.Getenv("PUBNUB_PUBLISH_KEY")
var SUBSCRIBE_KEY = os.Getenv("PUBNUB_SUBSCRIBE_KEY")
var SECRET_KEY = os.Getenv("PUBNUB_SECRET_KEY")
var CIPHER_KEY = ""
var USE_SSL = false
var PUBNUB_CHANNEL_ID = "hello_world"

func infinitePublish(pubnubChannelId string) {
	pubnub := messaging.NewPubnub(
		PUBLISH_KEY,
		SUBSCRIBE_KEY,
		SECRET_KEY,
		CIPHER_KEY,
		USE_SSL,
		"",  // custom UUID
		nil, // optional logger
	)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	go handlePublishCallbacks(successChannel, errorChannel)
	for {
		pubnub.Publish(
			pubnubChannelId,
			"Arbitrary payload, this can be anything",
			successChannel,
			errorChannel,
		)
		time.Sleep(2 * time.Second)
	}
}

func handlePublishCallbacks(successChannel, errorChannel chan []byte) {
	select {
	case response := <-successChannel:
		//fmt.Println(string(response))
	case err := <-errorChannel:
		fmt.Println(string(err))
	case <-messaging.Timeout():
		fmt.Println("Publish() timeout")
	}
}

func extractMessagesFromChannel(pubnubChannelId string, rawMessages chan string) {
	pubnub := messaging.NewPubnub(
		PUBLISH_KEY,
		SUBSCRIBE_KEY,
		SECRET_KEY,
		CIPHER_KEY,
		USE_SSL,
		"",  // custom UUID
		nil, // optional logger
	)
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)
	go pubnub.Subscribe(
		pubnubChannelId,
		"", // no idea
		successChannel,
		false,
		errorChannel,
	)
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

func main() {
	incomingMessages := make(chan string)
	go extractMessagesFromChannel(PUBNUB_CHANNEL_ID, incomingMessages)
	go infinitePublish(PUBNUB_CHANNEL_ID)
	for {
		payload := <-incomingMessages
		fmt.Printf("Got message: %s\n", payload)
	}

}
