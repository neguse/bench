
package main

import (
	"crypto/tls"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
	"sync"
	"fmt"
)

type TaskConfig struct {
	EndPoint string
}

type Result struct {
	e error
	d time.Duration
	rb []byte
}

func NewResult(e error, d time.Duration, rb []byte) *Result {
	r := new(Result)
	r.e = e
	r.d = d
	r.rb = rb
	return r
}

func task(wg *sync.WaitGroup, conf *TaskConfig, tr *http.Transport, res_ch chan *Result, quit_ch chan int) {
	defer wg.Done()
	client := &http.Client {
		Transport: tr,
	}
	body := bytes.NewBufferString(`{"jsonrpc":"2.0", "id":1,"method":"index","params":[1,2,3]}`)
	for {
		// resp.Body.Closeをdeferでやるための関数
		func () {
			t1 := time.Now()
			url := fmt.Sprintf("https://%s/index", conf.EndPoint)
			method := "application/json"
			resp, err := client.Post(url, method, body)
			if err != nil {
				res_ch <- NewResult(err, t1.Sub(time.Now()), make([]byte, 0))
				fmt.Println(err)
				continue;
			} else {
				defer resp.Body.Close()
				rb, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err, string(rb))
					res_ch <- NewResult(err, t1.Sub(time.Now()), rb)
					continue;
				}
				res_ch <- NewResult(nil, t1.Sub(time.Now()), rb)
			}
		}()

		select {
		// case <- time.After(time.Second * 1):
		//	continue
		case <- quit_ch:
			return
		default:

		}
	}
}

func receive(res_ch chan *Result) {
	count := 0
	dur := 0.0
	tick := time.Tick(time.Second * 1)
	done := false
	for (!done) {
		select {
		case res, more := <- res_ch:
			if (more) {
				count++
				dur -= res.d.Seconds()
			} else {
				done = true
			}
		case <- tick:
			if count != 0 {
				fmt.Println(count, dur / float64(count))
			}
			count = 0
			dur = 0.0
		}
	}
}

func main() {
  conf := TaskConfig { EndPoint: "10.0.0.125" }
	res_ch := make(chan *Result, 10000)
	quit_ch := make(chan int)
	tr := &http.Transport {
		DisableKeepAlives	: true,
		DisableCompression	: true,
		TLSClientConfig: &tls.Config {
			InsecureSkipVerify: true,
			CipherSuites: []uint16 {
				tls.TLS_RSA_WITH_AES_128_CBC_SHA
			},
			SessionTicketsDisabled: true
		},
		TLSHandshakeTimeout: time.Second * 30,
		ResponseHeaderTimeout: time.Second * 30,
	}
	go receive(res_ch)
	var wg sync.WaitGroup
	for i := 1; i < 10000; i++ {
		wg.Add(1)
		go task(&wg, &conf, tr, res_ch, quit_ch)
		if i % 100 == 0 {
			time.Sleep(time.Millisecond * 100)
			fmt.Println("launched ", i, " tasks")
		}
	}
	time.Sleep(time.Second * 100)
	fmt.Println("closing quit_ch")
	close(quit_ch);
	fmt.Println("waiting all tasks")
	wg.Wait()
	fmt.Println("done")
}

