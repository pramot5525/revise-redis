package main

import (
	"fmt"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func demoTTLAndEviction(rdb *goredis.Client) {
	fmt.Println("\n=== TTL & EVICTION ===")

	// --- TTL Basics ---
	fmt.Println("\n--- TTL Basics ---")

	// SET with TTL
	rdb.Set(ctx, "otp:user:42", "836492", 5*time.Minute)
	ttl, _ := rdb.TTL(ctx, "otp:user:42").Result()
	fmt.Printf("OTP expires in → %v\n", ttl.Round(time.Second))

	// key ไม่มี TTL → returns -1
	rdb.Set(ctx, "permanent:key", "value", 0)
	ttl2, _ := rdb.TTL(ctx, "permanent:key").Result()
	fmt.Printf("permanent key TTL → %v  (-1 = no expiry)\n", ttl2)

	// PERSIST — ลบ TTL ออก
	rdb.Set(ctx, "temp:key", "hello", 10*time.Second)
	rdb.Persist(ctx, "temp:key")
	ttl3, _ := rdb.TTL(ctx, "temp:key").Result()
	fmt.Printf("after PERSIST → %v  (-1 = no expiry)\n", ttl3)

	// EXPIREAT — หมดอายุตามเวลา (reset ตอนเที่ยงคืน)
	tomorrow := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	rdb.Set(ctx, "daily:counter", 0, 0)
	rdb.ExpireAt(ctx, "daily:counter", tomorrow)
	ttl4, _ := rdb.TTL(ctx, "daily:counter").Result()
	fmt.Printf("daily:counter resets in → %v\n", ttl4.Round(time.Hour))

	rdb.Del(ctx, "permanent:key", "temp:key", "daily:counter")

	// --- TTL Patterns ---
	fmt.Println("\n--- TTL Patterns ---")

	// Pattern 1: Short-lived token (reset password, OTP)
	rdb.Set(ctx, "tok:reset:99", "abc123", 15*time.Minute)
	t1, _ := rdb.TTL(ctx, "tok:reset:99").Result()
	fmt.Printf("reset token expires in %v\n", t1.Round(time.Second))

	// Pattern 2: Login attempt lockout
	attemptKey := "login:attempts:user:5"
	rdb.Del(ctx, attemptKey)
	for i := 1; i <= 5; i++ {
		cnt, _ := rdb.Incr(ctx, attemptKey).Result()
		if cnt == 1 {
			rdb.Expire(ctx, attemptKey, 15*time.Minute)
		}
		if cnt >= 5 {
			ttlLock, _ := rdb.TTL(ctx, attemptKey).Result()
			fmt.Printf("account locked after %d attempts! retry in %v\n", cnt, ttlLock.Round(time.Minute))
		}
	}

	// Pattern 3: Cache stampede prevention (SETNX as lock)
	lockKey := "lock:refresh:product:1"
	rdb.Del(ctx, lockKey)
	gotLock, _ := rdb.SetNX(ctx, lockKey, "1", 10*time.Second).Result()
	if gotLock {
		fmt.Println("got lock — refreshing cache from DB")
		rdb.Set(ctx, "product:1", `{"id":1,"name":"Mouse"}`, 5*time.Minute)
		rdb.Del(ctx, lockKey)
	} else {
		fmt.Println("another worker is refreshing, skip")
	}

	rdb.Del(ctx, "tok:reset:99", attemptKey, "product:1")

	// --- Memory Info ---
	fmt.Println("\n--- Memory & Eviction Info ---")
	info, _ := rdb.Info(ctx, "memory").Result()
	for _, line := range strings.Split(info, "\n") {
		for _, key := range []string{"used_memory_human", "used_memory_peak_human", "maxmemory_policy"} {
			if strings.HasPrefix(line, key+":") {
				fmt.Printf("  %s\n", strings.TrimSpace(line))
			}
		}
	}

	policy, _ := rdb.ConfigGet(ctx, "maxmemory-policy").Result()
	fmt.Println("  eviction policy:", policy)

	// Eviction Policies Reference:
	// noeviction   → คืน error เมื่อ RAM เต็ม (default)
	// allkeys-lru  → ลบ key ที่ไม่ได้ใช้นานที่สุด (ทุก key)
	// volatile-lru → ลบ key ที่มี TTL และไม่ได้ใช้นานที่สุด
	// allkeys-lfu  → ลบ key ที่ถูกใช้น้อยที่สุด
	// volatile-ttl → ลบ key ที่ TTL เหลือน้อยที่สุด
}
