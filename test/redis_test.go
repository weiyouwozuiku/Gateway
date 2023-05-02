package test

import (
	"fmt"
	server2 "github.com/weiyouwozuiku/Gateway/handler"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
)

func Test_Redis(t *testing.T) {
	SetUp()
	c, err := server2.RedisConnFactory("default")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	trace := public.NewTrace()
	redisKey := "test_key1"
	server.RedisConnDo(trace, c, "SET", redisKey, "hello")
	server.RedisConnDo(trace, c, "EXPIRE", "test_key1", 10)
	vint, verr := redis.Int64(server.RedisConnDo(trace, c, "INCR", "test_incr"))
	server.RedisConnDo(trace, c, "EXPIRE", "test_incr", 3600)
	fmt.Println(vint)
	if verr != nil {
		t.Fatal(verr)
	}
	v, err := redis.String(server.RedisConnDo(trace, c, "GET", redisKey))
	fmt.Println(v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello" || err != nil {
		t.Fatal("test redis get fatal!")
	}
	TearDown()
}
