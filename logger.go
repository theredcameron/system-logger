package systemlogger

import(
	"fmt"
	"os"
	"strings"
	"time"
	"log"
)

type Logger struct {
	LogFilePath	string
	LogFile		*os.File
	LogFileName string
	EnableDebug bool
}

func CreateLogger(filePathInput string, debug bool) (*Logger, error) {
	filePath := strings.TrimSpace(filePathInput)

	if len(filePath) == 0 {
		return nil, fmt.Errorf("The log file path is required")
	}

	return &Logger{
		LogFilePath: filePath,
		EnableDebug: debug,
	}, nil
}

func (this *Logger) Log(message string) {
	message = fmt.Sprintf("INFO: %v", message)
	this.log(message)
}

func (this *Logger) ErrorMessage(message string) {
	message = fmt.Sprintf("ERROR: %v", message)
	this.log(message)
}

func (this *Logger) Error(err error) {
	this.ErrorMessage(err.Error())
}

func (this *Logger) Debug(message string) {
	if this.EnableDebug {
		message = fmt.Sprintf("DEBUG: %v", message)
		this.log(message)
	}
}

func (this *Logger) log(message string) {
	logFileName := time.Now().Format("20060102") + ".log"

	_, err := os.ReadDir(this.LogFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(this.LogFilePath, 0750)
			if err != nil {
				panic(err)
			}
		}
	}

	if logFileName != this.LogFileName {
		this.LogFile.Close()
		newFile, err := os.OpenFile(this.LogFilePath + "/" + logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		this.LogFileName = logFileName
		this.LogFile = newFile

		log.SetOutput(this.LogFile)
	}
	
	log.Println(message)
}