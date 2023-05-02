package test

import (
	"fmt"
	server2 "github.com/weiyouwozuiku/Gateway/handler"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/handler"
	"github.com/weiyouwozuiku/Gateway/public"
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
	handler.RedisConnDoWithCtx(trace, c, "SET", redisKey, "hello")
	handler.RedisConnDoWithCtx(trace, c, "EXPIRE", "test_key1", 10)
	vint, verr := redis.Int64(handler.RedisConnDoWithCtx(trace, c, "INCR", "test_incr"))
	handler.RedisConnDoWithCtx(trace, c, "EXPIRE", "test_incr", 3600)
	fmt.Println(vint)
	if verr != nil {
		t.Fatal(verr)
	}
	v, err := redis.String(handler.RedisConnDoWithCtx(trace, c, "GET", redisKey))
	fmt.Println(v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello" || err != nil {
		t.Fatal("test redis get fatal!")
	}
	TearDown()
}
