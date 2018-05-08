package main

import (
	"strings"
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
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
	for v := range wc {
		fmt.Println(v)
	}
}

func (r *ReadFromFile) Read(rc chan string) {
	//Read Module

	//open file
	file, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}

	//从文件末尾逐行开始读取文件内容
	file.Seek(0,2)
	rd := bufio.NewReader(file)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}
		rc <- string(line)
	}

}

type LogProcess struct {
	rc chan string
	wc chan string
	read Reader
	write Writer
}

func (l *LogProcess) Process(){
	//process
	for v := range l.rc {
		l.wc <- strings.ToUpper(v)
	}
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
