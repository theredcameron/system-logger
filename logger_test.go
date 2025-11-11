package systemlogger

import (
	"testing"
        "time"
)

func TestHelloWorld(t *testing.T) {
	loggerConfig := LoggerConfig{
            LogDirectory: "./logs",
            Debug: true,
            MaxTimespanInDays: 1,
        }

        logger, err := CreateLogger(loggerConfig)
	if err != nil {
		t.Fatal(err)
	}

        time.Sleep(5 * time.Minute)

	logger.Log("Logging Tests Complete")
}
