package middleware

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/weiyouwozuiku/Gateway/log"
	"github.com/weiyouwozuiku/Gateway/public"
	"github.com/weiyouwozuiku/Gateway/server"
	"sync/atomic"
	"time"
)

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
	return redis.Int64(server.RedisConfDo("GET", s.GetDayKey(t)))
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
			dayKey := reqCounter.Get
		}
	}()
	return service
}
