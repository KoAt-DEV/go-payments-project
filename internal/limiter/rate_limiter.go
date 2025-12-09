package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client redis.UniversalClient
}

var gcraScript = redis.NewScript(`
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local period = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	local tat = redis.call("GET", key)
	if not tat then
		tat = 0
	else
		tat = tonumber(tat)
	end

	local new_tat = math.max(tat, now) + period

	if new_tat - now > period * limit then
		local delay_ns = new_tat - now
		local delay_ms = math.floor(delay_ns / 1000000)
		return -delay_ms
	end

	redis.call("SET", key, new_tat, "PX", math.max(1, period * 2 / 1000000))
	return 0
`)

func New(client redis.UniversalClient) *RateLimiter {
	return &RateLimiter{client: client}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string, burst int, ratePerSec float64) (bool, time.Duration) {
	if ratePerSec <= 0 {
		return false, 0
	}

	periodNS := int64(1_000_000_000 / ratePerSec)
	nowNS := time.Now().UnixNano()

	result, err := gcraScript.Run(ctx, rl.client, []string{"rl:" + key}, burst, periodNS, nowNS).Int64()
	if err != nil && err != redis.Nil {
		return true, 0
	}

	if result < 0 {
		return false, time.Duration(-result) * time.Millisecond
	}

	return true, 0
}
