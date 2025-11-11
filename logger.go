package systemlogger

import(
	"fmt"
	"strings"
	"time"
	"log"
        "context"
        "os"
        "os/signal"
        "path/filepath"
)

type Logger struct {
	LogFilePath	string
	LogFile		*os.File
	LogFileName string
	EnableDebug bool
        MaxTimespanInDays int
}

func CreateLogger(filePathInput string, debug bool, timeInSeconds int) (*Logger, error) {
	filePath := strings.TrimSpace(filePathInput)

	if len(filePath) == 0 {
		return nil, fmt.Errorf("The log file path is required")
	}

	logger := &Logger{
		LogFilePath: filePath,
		EnableDebug: debug,
	}

        go logger.startLogCleaner(timeInSeconds)

        return logger, nil
}

func (this *Logger) startLogCleaner(timeInSeconds int) {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    ticker := time.NewTicker(time.Duration(timeInSeconds) * time.Second)
    defer ticker.Stop()

    this.Log("Logging Cleanup Started")

    for {
        select {
        case <-ticker.C:
            this.Log("cleanup")
        case <-ctx.Done():
            this.Log("Logging Cleanup Closing")
            return
        }
    }
}

func (this *Logger) logCleaningAction() error {
    err := filepath.WalkDir(this.LogFilePath, d fs.DirEntry, err error) {
        if err != nil {
            return err
        }

        if !d.IsDir() && d.Name() != this.LogFileName {
            info, err := d.Info()
            if err != nil {
                return err
            }

            modTime := info.ModTime()
            oldestDate := time.Now().Add((-1 * this.MaxTimespanInDays), * time.Day)
            if modTime.Before(oldestDate) {
                err = this.deleteFile(d.Name())
                if err != nil {
                    return err
                }
            }
        }

        return nil
    })

    if err != nil {
        return err
    }

    return nil
}

func (this *Logger) deleteFile(fileName string) error {
    err := os.Remove(fmt.Sprintf("%s/%s", this.LogFilePath, fileName))
    if != err {
        return err
    }

    return nil
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
