package app

import (
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors        map[string]*visitor
	mu              sync.RWMutex
	rateLimitPerSec int
	rateBurst       int
	Middleware      func(next http.Handler) http.Handler
}

func NewRateLimiter(rateLimitPerSec int, rateBurst int) *RateLimiter {
	rl := &RateLimiter{
		visitors:        make(map[string]*visitor),
		mu:              sync.RWMutex{},
		rateBurst:       rateBurst,
		rateLimitPerSec: rateLimitPerSec,
	}
	rl.Middleware = func(next http.Handler) http.Handler {
		return rl.limiterMiddleware(next)
	}

	go rl.cleanupVisitors()
	return rl
}

const (
	ttl         = 5 * time.Minute
	cleanupCron = 2 * time.Minute
)

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(rl.rateLimitPerSec), rl.rateBurst)
		rl.visitors[ip] = &visitor{limiter, common.Now()}
		return limiter
	}

	v.lastSeen = common.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(cleanupCron)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > ttl {
				log.Debug().Str("ip", ip).Msg("Cleaning up rate limiter visitor")
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// with the help of https://www.alexedwards.net/blog/how-to-rate-limit-http-requests, TY!
func (rl *RateLimiter) limiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ipFrom(r)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			log.Error().Str("ip", ip).Msg("Someone spamming? Rate limit hit!")
			rest.TooManyRequests(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
