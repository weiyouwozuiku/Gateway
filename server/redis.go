package server

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/public"
)

var (
	RedisConfMap *RedisMapConf
)

type RedisConf struct {
	ProxyList    []string `mapstructure:"proxy_list"`
	Password     string   `mapstructure:"password"`
	DB           int      `mapstructure:"db"`
	ConnTimeout  int      `mapstructure:"conn_timeout"`
	ReadTimeout  int      `mapstructure:"read_timeout"`
	WriteTimeout int      `mapstructure:"write_timeout"`
}
type RedisMapConf struct {
	List map[string]*RedisConf `mapstructure:"list"`
}

func InitRedisConf(confName string) error {
	RedisConfMap = &RedisMapConf{}
	if err := public.ParseConfig(confName, RedisConfMap); err != nil {
		return err
	}
	return nil
}

func RedisConnFactory(name string) (redis.Conn, error) {
	if RedisConfMap != nil && RedisConfMap.List != nil {
		for confName, cfg := range RedisConfMap.List {
			if name == confName {
				randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
				// redis配置缺失时，设置默认参数
				if cfg.ConnTimeout == 0 {
					cfg.ConnTimeout = 50
				}
				if cfg.ReadTimeout == 0 {
					cfg.ReadTimeout = 100
				}
				if cfg.WriteTimeout == 0 {
					cfg.WriteTimeout = 100
				}
				c, err := redis.Dial(
					"tcp",
					randHost,
					redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
				)
				if err != nil {
					return nil, err
				}
				if cfg.Password != "" {
					if _, err := c.Do("AUTH", cfg.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
				if cfg.DB != 0 {
					if _, err := c.Do("SELECT", cfg.DB); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, nil
			}
		}
	}
	return nil, errors.New("create redis conn fail")
}
func RedisConnDo(trace *public.TraceContext, c redis.Conn, commandName string, args ...any) (any, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		public.Log.TagError(trace, public.LTagRedisFailed, map[string]any{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		public.Log.TagInfo(trace, public.LTagRedisSuccess, map[string]any{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}

// 通过配置 执行redis
func RedisConfDo(trace *public.TraceContext, name string, commandName string, args ...any) (any, error) {
	c, err := RedisConnFactory(name)
	defer c.Close()
	if err != nil {
		public.Log.TagError(trace, public.LTagRedisFailed, map[string]any{
			"method": commandName,
			"err":    errors.New("redisConfFactory_error:" + name),
			"bind":   args,
		})
		return nil, err
	}
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		public.Log.TagError(trace, public.LTagRedisFailed, map[string]any{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		public.Log.TagInfo(trace, public.LTagRedisSuccess, map[string]any{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}

func RedisLogDo(trace *public.TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		public.TagError(trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}
