package middleware

import (
	"fmt"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimitMiddleware(next http.HandlerFunc, limit *ratelimit.Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 尝试获取一个请求令牌
		limit.Wait(1)

		// 计算处理时间
		elapsedTime := time.Since(startTime)
		fmt.Printf("Processing time: %dms", elapsedTime.Milliseconds())

		// 调用下一个处理器
		next(w, r)
	}
}
