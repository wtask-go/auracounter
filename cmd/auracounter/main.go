package main

import (
	"log"
	"os"
)

func main() {

	l := log.New(os.Stdout, "rpcsrv ", log.LUTC|log.Ldate|log.Lmicroseconds|log.Ltime|log.Llongfile)
	l.Output(2, "Hello!")
}
