package middleware

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/handler"
	"github.com/weiyouwozuiku/Gateway/log"
	"github.com/weiyouwozuiku/Gateway/public"
	"sync/atomic"
	"time"
)

const ONEDAYSECOND = 60 * 60 * 24

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func (s *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(public.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", public.RedisFlowDayKey, dayStr, s.AppID)
}
func (s *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(public.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", public.RedisFlowHourKey, hourStr, s.AppID)
}
func (s *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(handler.RedisConfDo("GET", s.GetDayKey(t)))
}
func (s *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(handler.RedisConfDo("GET", s.GetDayKey(t)))
}
func (s *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("RedisFlowCountService Increase error||error=%v", err)
			}
		}()
		atomic.AddInt64(&s.TickerCount, 1)
	}()
}
func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("NewRedisFlowCountService error||error=%v", err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) // 获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据
			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)
			if err := handler.RedisConfPipline(func(c redis.Conn) {
				c.Send("INCRBY", dayKey, tickerCount)
				c.Send("EXPIRE", dayKey, ONEDAYSECOND*2)
				c.Send("INCRBY", hourKey, tickerCount)
				c.Send("EXPIRE", hourKey, ONEDAYSECOND*2)
			}); err != nil {
				log.Error("PipRedis of NewRedisFlowCountService error||error=%v", err)
			}
			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				log.Error("reqCounter.GetDayData err||error=%v", err)
				continue
			}
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = nowUnix
				continue
			}
			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = nowUnix
			}
		}
	}()
	return reqCounter
}
