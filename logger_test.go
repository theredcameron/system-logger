package systemlogger

import (
	"testing"
        "time"
)

func TestHelloWorld(t *testing.T) {
	logger, err := CreateLogger("./logs", true, 2)
	if err != nil {
		t.Fatal(err)
	}

        time.Sleep(10 * time.Second)

	logger.Log("test log")
}
