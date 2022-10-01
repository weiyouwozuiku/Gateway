package test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/weiyouwozuiku/Gateway/initialize"
)

var (
	addr       string    = "127.0.0.1"
	initOnce   sync.Once = sync.Once{}
	serverOnce sync.Once = sync.Once{}
)

// 初始化测试用例
func SetUp() {
	initOnce.Do(func() {
		if err := initialize.InitModules("../conf/dev/"); err != nil {
			log.Fatal(err)
		}
	})
}

func TearDown() {
	initialize.Destory()
}

// 测试Post请求
func TestPost(t *testing.T) {
	InitTestServer()
}

func InitTestServer() {
	serverOnce.Do(func() {
		http.HandleFunc("/postjson", func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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
	if err := initialize.InitModules("../conf/dev/"); err != nil {
		fmt.Println(err)
	}
}
