package main

import (
	"strings"
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"regexp"
	"log"
	"strconv"
	"net/url"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan *Message)
}

type ReadFromFile struct {
	path string //read file path
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

type Message struct {
	TimeLocal						time.Time
	BytesSent						int
	Path, Method, Scheme, Status	string
	UpstreamTime, RequestTime		float64
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
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
		rc <- string(line[:len(line)-1])
	}

}

type LogProcess struct {
	rc chan string
	wc chan *Message
	read Reader
	write Writer
}

func (l *LogProcess) Process(){
	//process
	//从Read Channel中读取每行日志数据
	//正则来对其进行匹配
	/**
	172.0.0.12 -- [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-"1.005 1.854
	
	([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	**/

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	location, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(v)
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", v)
			continue
		}
		
		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], location)

		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		//GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail", ret[6])
			continue
		}
		message.Method = reqSli[0]
		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err)
			continue
		}
		message.Path = u.Path

		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		l.wc <- message
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
		wc: make(chan *Message),
		read: r,
		write: w,
	}

	go lp.read.Read(lp.rc)

	go lp.Process()

	go lp.write.Write(lp.wc)

	time.Sleep(30*time.Second)

}
