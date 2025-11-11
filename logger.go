package systemlogger

import(
	"fmt"
	"strings"
	"time"
	"log"
        "context"
        "io/fs"
        "os"
        "os/signal"
        "path/filepath"
)

type LoggerConfig struct {
    LogDirectory         string      `json:"logDirectory"`
    Debug                bool        `json:"debug"`
    MaxTimespanInDays    int         `json:"maxTimespanInDays"`
}

type Logger struct {
	LogFilePath	    string
	LogFile		    *os.File
	LogFileName         string
	EnableDebug         bool
        MaxTimespanInDays   int
}

func CreateLogger(loggerConfig LoggerConfig) (*Logger, error) {
	filePath := strings.TrimSpace(loggerConfig.LogDirectory)

	if len(filePath) == 0 {
		return nil, fmt.Errorf("The log file path is required")
	}

	logger := &Logger{
		LogFilePath: filePath,
		EnableDebug: loggerConfig.Debug,
                MaxTimespanInDays: loggerConfig.MaxTimespanInDays,
	}

        go logger.startLogCleaner()

        return logger, nil
}

func (this *Logger) startLogCleaner() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    ticker := time.NewTicker(time.Duration(6) * time.Hour) //Check the file directory every six hours
    defer ticker.Stop()

    this.metaLog("Logging Cleanup Started")

    for {
        select {
        case <-ticker.C:
            err := this.logCleaningAction()
            if err != nil {
                this.metaLog(fmt.Sprintf("ERROR: %v", err))
            }
        case <-ctx.Done():
            this.metaLog("Logging Cleanup Closing")
            return
        }
    }
}

func (this *Logger) logCleaningAction() error {
    err := filepath.WalkDir(this.LogFilePath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if !d.IsDir() && d.Name() != this.LogFileName {
            info, err := d.Info()
            if err != nil {
                return err
            }

            modTime := info.ModTime()
            oldestDate := time.Now().Add(time.Duration((-1 * this.MaxTimespanInDays)) * (24 * time.Hour)) //Set the max age in days for the log files
            if modTime.Before(oldestDate) {
                err = this.deleteFile(d.Name())
                if err != nil {
                    return err
                }

                this.metaLog(fmt.Sprintf("Deleted log file: %s", d.Name()))
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
    if err != nil {
        return err
    }

    return nil
}

func (this *Logger) metaLog(message string) {
    message = fmt.Sprintf("***LOGGING SYSTEM***: %s", message)
    this.log(message)
}

func (this *Logger) Log(message string) {
	message = fmt.Sprintf("INFO: %s", message)
	this.log(message)
}

func (this *Logger) ErrorMessage(message string) {
	message = fmt.Sprintf("ERROR: %s", message)
	this.log(message)
}

func (this *Logger) Error(err error) {
	this.ErrorMessage(err.Error())
}

func (this *Logger) Debug(message string) {
	if this.EnableDebug {
		message = fmt.Sprintf("DEBUG: %s", message)
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
