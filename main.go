package main

import (
	"fmt"
)

func main() {
	incomingMessages := make(chan string)
	go extractMessagesFromChannel(PUBNUB_CHANNEL_ID, incomingMessages)
	go infinitePublish(PUBNUB_CHANNEL_ID)
	for {
		payload := <-incomingMessages
		fmt.Printf("Got message: %s\n", payload)
	}

}
