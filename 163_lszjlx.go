//批量获取雅虎股票数据。
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	UA = "Golang Downloader from Ijibu.com"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	logfile, _ := os.OpenFile("./log/163_lszjlx_sh.log", os.O_RDWR|os.O_CREATE, 0)
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	c := make(chan int, 941)
	fh, ferr := os.Open("./shang_new.ini")
	if ferr != nil {
		return
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)

	for i := 1; i <= 942; i++ { //加入goroutine缓冲，4个执行完了再执行下面的4个
		input, _ := inputread.ReadString('\n')
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "0")
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "50")
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "100")
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "150")
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "200")
		go func(logger *log.Logger, logfile *os.File, input string, date string) {
			getShangTickerTables(logger, logfile, input, date)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input), "250")

		time.Sleep(6 * time.Second) //加入执行缓冲，否则同时发起大量的tcp连接，操作系统会直接返回错误。

	}
	defer logfile.Close()
	for j := 0; j < 941; j++ {
		<-c
	}
}

func getShangTickerTables(logger *log.Logger, logfile *os.File, code string, date string) {
	//并发写文件必须要有锁啊，怎么还是串行程序的思维啊。
	fileName := "./data/163/lszjlx/sh/" + code + "_" + date + ".html"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()

	urls := "http://quotes.money.163.com/trade/lszjlx_" + code + "," + date + ".html"
	var req http.Request
	req.Method = "GET"
	req.Close = true
	req.URL, err = url.Parse(urls)
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("User-Agent", UA)
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err == nil {
		if resp.StatusCode == 200 {
			logger.Println(logfile, code+":sucess"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":sucess")
			io.Copy(f, resp.Body)
		} else {
			logger.Println(logfile, code+":http get StatusCode"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":" + strconv.Itoa(resp.StatusCode))
		}
		defer resp.Body.Close()
	} else {
		logger.Println(logfile, code+":http get error"+code)
		fmt.Println(code + ":error")
	}
}