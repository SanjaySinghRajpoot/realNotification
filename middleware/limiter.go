package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
)

const (
	maxRequests     = 30
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

	// get from redis cache
	count, err := utils.GetIPAddress(ip)

	if err != nil {
		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache: %d", count)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	if count >= maxRequests {
		context.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	count = count + 1

	msg, err := utils.SetIPAddress(ip, count)
	if err != nil {
		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache: %s", msg)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		count, err := utils.GetIPAddress(ip)

		if err != nil {
			if err != nil {
				fmt.Printf("Failed to Get the Redis Cache: %d", count)

				context.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
		}

		count = count - 1

		msg, err := utils.SetIPAddress(ip, count)
		if err != nil {
			if err != nil {
				fmt.Printf("Failed to Set the Redis Cache: %s", msg)

				context.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
		}

	})

	context.Next()
}
