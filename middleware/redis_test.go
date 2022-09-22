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
	middleware.RedisConnDo(trace, c, "SET", redisKey, "test_redis")
	middleware.RedisConnDo(trace, c, "EXPIRE", "test_key1", 10)
	vint, verr := redis.Int64(middleware.RedisConnDo(trace, c, "INCR", "test_incr"))
	fmt.Println(vint)
	if verr != nil {
		t.Fatal(verr)
	}
	v, err := redis.String(middleware.RedisConnDo(trace, c, "GET", redisKey))
	fmt.Println(v)
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
	if v != "test_redis" || err != nil {
		t.Fatal("test redis get fatal!")
	}
	TearDown()
}
