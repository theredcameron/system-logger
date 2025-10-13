package systemlogger

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	logger, err := CreateLogger("./logs", true)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log("test log")
}