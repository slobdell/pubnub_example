package main

import (
	"fmt"
	"github.com/pubnub/go/messaging"
	"time"
)

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
	case <-successChannel:
		//fmt.Println(string(response))
	case err := <-errorChannel:
		fmt.Println(string(err))
	case <-messaging.Timeout():
		fmt.Println("Publish() timeout")
	}
}
