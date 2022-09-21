package middleware_test

import (
	"bytes"
	"fmt"
	"gateway/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var (
	addr       string    = "127.0.0.1"
	initOnce   sync.Once = sync.Once{}
	serverOnce sync.Once = sync.Once{}
)

// 初始化测试用例
func SetUp() {
	initOnce.Do(func() {
		if err := middleware.InitModules("../conf/dev/"); err != nil {
			log.Fatal(err)
		}
	})
}

func TearDown() {
	middleware.Destory()
}

// 测试Post请求
func TestPost(t *testing.T) {
	InitTestServer()
}

func InitTestServer() {
	serverOnce.Do(func() {
		http.HandleFunc("/postjson", func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			w.Write([]byte(bodyBytes))
		})
		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			cityId := r.FormValue("city_id")
			w.Write([]byte(cityId))
		})
		http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			cityId := r.FormValue("city_id")
			w.Write([]byte(cityId))
		})
		go func() {
			log.Println("ListenAndServer ", addr)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				log.Fatal("ListenAndServer: ", err)
			}
		}()
		time.Sleep(time.Second)
	})
}

func Test_InitModules(t *testing.T) {
	// pwd, _ := os.Getwd()
	// fmt.Println(pwd)
	if err := middleware.InitModules("../conf/dev/"); err != nil {
		fmt.Println(err)
	}
}
