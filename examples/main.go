package main

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("connect failed: %v", err)
	}

	demoStrings(rdb)
	demoLists(rdb)
	demoSets(rdb)
	demoSortedSets(rdb)
	demoHashes(rdb)
	demoCaching(rdb)
	demoSessionManagement(rdb)
	demoRateLimiting(rdb)
	demoPubSub(rdb)
	demoTTLAndEviction(rdb)
}
