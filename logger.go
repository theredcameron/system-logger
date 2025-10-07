package systemlogger

import(
	"fmt"
)

type Logger struct {

}

func (this *Logger) HelloWorld() {
	fmt.Println("Hello world")
}