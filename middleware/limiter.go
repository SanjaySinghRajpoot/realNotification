package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxRequests     = 20
	perMinutePeriod = 1 * time.Minute
)

var (
	// update this with DB query
	ipRequestsCounts = make(map[string]int)
	mutex            = &sync.Mutex{}
)

func RateLimiter(context *gin.Context) {
	ip := context.ClientIP()
	mutex.Lock()
	defer mutex.Unlock()
	count := ipRequestsCounts[ip]
	if count >= maxRequests {
		context.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	ipRequestsCounts[ip] = count + 1
	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
	})

	context.Next()
}
