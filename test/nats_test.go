package test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

const NATS_PORT = "nats://localhost:4222"

func Test_Nats(t *testing.T) {
	nc, err := nats.Connect(NATS_PORT)
	if err != nil {
		t.Fatalf("Failed to connect to nats server: %v", err)
	}
	defer nc.Close()

	// Subscribe to a subject
	msgChannel := make(chan *nats.Msg)

	sub, err := nc.Subscribe("test", func(msg *nats.Msg) {
		msgChannel <- msg
	})

	if err != nil {
		t.Fatalf("Failed to subscribe to subject: %v", err)
	}

	// Publish a message
	err = nc.Publish("test", []byte("Message Nikola!"))
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for message to be received or for timeout to occur
	timeout := time.After(5 * time.Second)
	var receivedMsg *nats.Msg

	select {
	case receivedMsg = <-msgChannel:
		// Message received, continue with test
	case <-timeout:
		sub.Unsubscribe()
		t.Fatal("Timeout waiting for message to be received")
	}

	// Check that the message is correct
	if string(receivedMsg.Data) != "Message Nikola!" {
		t.Errorf("Received unexpected message: %v", receivedMsg.Data)
	}

	// Unsubscribe
	sub.Unsubscribe()

}
