package main

import (
	"strings"
	"fmt"
	"time"
)

type LogProcess struct {
	rc chan string
	wc chan string
	path string //file path
	influxDBDsn string // influx data source
}

func (l *LogProcess) ReadFromFile(){
	//read module
	l.rc <- "Hello, Dylan"
}

func (l *LogProcess) Process(){
	//process
	ct := <- l.rc
	l.wc <- strings.ToUpper(ct)
}

func (l *LogProcess) WriteToInfluxDB(){
	//write module
	fmt.Println(<-l.wc)
}

func main() {
	lp := &LogProcess{
		rc: make(chan string),
		wc: make(chan string),
		path: ".\\access.log",
		influxDBDsn: "username&password..",
	}

	go lp.ReadFromFile()

	go lp.Process()

	go lp.WriteToInfluxDB()

	time.Sleep(1*time.Second)

}
