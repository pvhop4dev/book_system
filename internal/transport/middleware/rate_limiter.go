package middleware

import (
	"book_system/internal/baselib/concurrentmap"
	"book_system/internal/config"
	"book_system/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		rl, ok := mapRateLimiter.Get(c.ClientIP())
		if !ok {
			cfg := config.MustGet()
			rl = &rateLimiter{
				limiter: rate.NewLimiter(rate.Limit(cfg.RateLimiter.Rate), cfg.RateLimiter.Burst),
			}
			mapRateLimiter.Set(c.ClientIP(), rl)
		}
		rl.lastRequest = time.Now()

		if !rl.limiter.AllowN(time.Now(), 1) {
			err := TooManyRequests
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    err.Code,
				"message": err.GetMesssageI18n(utils.GetCurrentLang(c)),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

var mapRateLimiter concurrentmap.ConcurrentMap[string, *rateLimiter]

type rateLimiter struct {
	limiter     *rate.Limiter
	lastRequest time.Time
}

func init() {
	mapRateLimiter = concurrentmap.New[*rateLimiter]()
	go func() {
		cleanLimiter()
	}()
}

func cleanLimiter() {
	for {
		time.Sleep(time.Second * 30)
		for k, v := range mapRateLimiter.Items() {
			if time.Since(v.lastRequest) > time.Minute*5 {
				mapRateLimiter.Remove(k)
			}
		}
	}
}
