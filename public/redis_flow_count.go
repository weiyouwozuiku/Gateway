package public

import "time"

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	service := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {

			}
		}()
	}()
	return service
}
