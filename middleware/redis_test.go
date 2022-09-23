package middleware_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/middleware"
)

func Test_Redis(t *testing.T) {
	SetUp()
	c, err := middleware.RedisConnFactory("default")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	trace := middleware.NewTrace()
	redisKey := "test_key1"
	middleware.RedisConnDo(trace, c, "SET", redisKey, "hello")
	middleware.RedisConnDo(trace, c, "EXPIRE", "test_key1", 10)
	vint, verr := redis.Int64(middleware.RedisConnDo(trace, c, "INCR", "test_incr"))
	middleware.RedisConnDo(trace, c, "EXPIRE", "test_incr", 3600)
	fmt.Println(vint)
	if verr != nil {
		t.Fatal(verr)
	}
	v, err := redis.String(middleware.RedisConnDo(trace, c, "GET", redisKey))
	fmt.Println(v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello" || err != nil {
		t.Fatal("test redis get fatal!")
	}
	TearDown()
}
