package systemlogger

import (
	"testing"
        "time"
)

func TestLog(t *testing.T) {
	loggerConfig := LoggerConfig{
            LogDirectory: "./logs",
            Debug: true,
            MaxTimespanInDays: 1,
        }

        logger, err := CreateLogger(loggerConfig)
	if err != nil {
		t.Fatal(err)
	}

        time.Sleep(5 * time.Second)

	logger.Log("Logging message")
	logger.Debug("Logging debug")
        logger.ErrorMessage("Logging error")
}
