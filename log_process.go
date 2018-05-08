package main

import (
	"strings"
	"fmt"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

type ReadFromFile struct {
	path string //read file path
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//write module
	fmt.Println(<-wc)
}

func (r *ReadFromFile) Read(rc chan string) {
	//Read Module
	rc <- "Hello, Dylan"
}

type LogProcess struct {
	rc chan string
	wc chan string
	read Reader
	write Writer
}

func (l *LogProcess) Process(){
	//process
	ct := <- l.rc
	l.wc <- strings.ToUpper(ct)
}

func main() {
	r := &ReadFromFile{
		path: ".\\access.log",
	}

	w := &WriteToInfluxDB{
		influxDBDsn: "username&password..",
	}

	lp := &LogProcess{
		rc: make(chan string),
		wc: make(chan string),
		read: r,
		write: w,
	}

	go lp.read.Read(lp.rc)

	go lp.Process()

	go lp.write.Write(lp.wc)

	time.Sleep(1*time.Second)

}
