package main

import (
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func demoCaching(rdb *goredis.Client) {
	fmt.Println("\n=== CACHING (Cache-Aside) ===")
	cacheKey := "product:42"
	rdb.Del(ctx, cacheKey)

	// รอบแรก — cache miss
	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == goredis.Nil {
		data := `{"id":42,"name":"Keyboard","price":1299}`
		rdb.Set(ctx, cacheKey, data, 5*time.Minute)
		fmt.Println("cache MISS — fetched from DB:", data)
	} else {
		fmt.Println("cache HIT:", cached)
	}

	// รอบสอง — cache hit
	cached, _ = rdb.Get(ctx, cacheKey).Result()
	fmt.Println("second read (HIT):", cached)

	rdb.Del(ctx, cacheKey)
	fmt.Println("cache invalidated")
}

func demoSessionManagement(rdb *goredis.Client) {
	fmt.Println("\n=== SESSION MANAGEMENT ===")
	sessionID := "sess:abc123"
	rdb.Del(ctx, sessionID)

	rdb.HSet(ctx, sessionID,
		"user_id", 42,
		"username", "alice",
		"role", "admin",
		"logged_in_at", time.Now().Unix(),
	)
	rdb.Expire(ctx, sessionID, 30*time.Minute)

	userID, _ := rdb.HGet(ctx, sessionID, "user_id").Result()
	ttl, _ := rdb.TTL(ctx, sessionID).Result()
	fmt.Printf("user_id: %s, expires in: %v\n", userID, ttl.Round(time.Second))

	// Sliding expiration — รีเซ็ต TTL ทุกครั้งที่ใช้งาน
	rdb.Expire(ctx, sessionID, 30*time.Minute)
	fmt.Println("TTL extended (sliding expiration)")

	rdb.Del(ctx, sessionID)
	fmt.Println("logged out — session deleted")
}

func demoRateLimiting(rdb *goredis.Client) {
	fmt.Println("\n=== RATE LIMITING (Fixed Window) ===")
	userKey := "user:99"
	limit := 5
	window := time.Minute

	for i := 1; i <= 7; i++ {
		allowed, remaining := checkRateLimit(rdb, userKey, limit, window)
		if allowed {
			fmt.Printf("request #%d → ✅ allowed  (remaining: %d)\n", i, remaining)
		} else {
			fmt.Printf("request #%d → ❌ rate limited\n", i)
		}
	}
	rdb.Del(ctx, fmt.Sprintf("ratelimit:%s", userKey))
}

func checkRateLimit(rdb *goredis.Client, key string, limit int, window time.Duration) (bool, int) {
	k := fmt.Sprintf("ratelimit:%s", key)
	count, _ := rdb.Incr(ctx, k).Result()
	if count == 1 {
		rdb.Expire(ctx, k, window)
	}
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return int(count) <= limit, remaining
}

func demoPubSub(rdb *goredis.Client) {
	fmt.Println("\n=== PUB/SUB ===")
	channel := "notifications"

	sub := rdb.Subscribe(ctx, channel)
	received := make(chan string, 5)
	go func() {
		for msg := range sub.Channel() {
			received <- msg.Payload
		}
	}()

	time.Sleep(50 * time.Millisecond)

	events := []string{"order:created:1001", "order:paid:1001", "order:shipped:1001"}
	for _, e := range events {
		rdb.Publish(ctx, channel, e)
		fmt.Println("published →", e)
	}

	time.Sleep(100 * time.Millisecond)
	sub.Close()
	close(received)

	fmt.Print("received:  ")
	for msg := range received {
		fmt.Printf("[%s] ", msg)
	}
	fmt.Println()
}
