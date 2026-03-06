package main

import (
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

func demoStrings(rdb *goredis.Client) {
	fmt.Println("\n=== STRINGS ===")

	rdb.Set(ctx, "user:name", "Alice", 0)
	val, _ := rdb.Get(ctx, "user:name").Result()
	fmt.Println("GET user:name →", val)

	rdb.Set(ctx, "page:views", 0, 0)
	rdb.Incr(ctx, "page:views")
	rdb.IncrBy(ctx, "page:views", 9)
	count, _ := rdb.Get(ctx, "page:views").Int()
	fmt.Println("page:views →", count)

	acquired, _ := rdb.SetNX(ctx, "lock:job", "worker-1", 0).Result()
	fmt.Println("acquired lock →", acquired)
	again, _ := rdb.SetNX(ctx, "lock:job", "worker-2", 0).Result()
	fmt.Println("second attempt →", again)
	rdb.Del(ctx, "lock:job")

	rdb.MSet(ctx, "k1", "v1", "k2", "v2", "k3", "v3")
	vals, _ := rdb.MGet(ctx, "k1", "k2", "k3").Result()
	fmt.Println("MGET →", vals)
}

func demoLists(rdb *goredis.Client) {
	fmt.Println("\n=== LISTS ===")
	rdb.Del(ctx, "queue:tasks", "stack:pages", "feed")

	rdb.RPush(ctx, "queue:tasks", "task1", "task2", "task3")
	task, _ := rdb.LPop(ctx, "queue:tasks").Result()
	fmt.Println("dequeue (FIFO) →", task)

	rdb.LPush(ctx, "stack:pages", "page1", "page2", "page3")
	top, _ := rdb.LPop(ctx, "stack:pages").Result()
	fmt.Println("stack pop (LIFO) →", top)

	rdb.RPush(ctx, "feed", "post:1", "post:2", "post:3", "post:4", "post:5")
	feed, _ := rdb.LRange(ctx, "feed", 0, 2).Result()
	fmt.Println("feed top 3 →", feed)

	length, _ := rdb.LLen(ctx, "feed").Result()
	fmt.Println("feed length →", length)
}

func demoSets(rdb *goredis.Client) {
	fmt.Println("\n=== SETS ===")
	rdb.Del(ctx, "tags:post:1", "tags:post:2")

	rdb.SAdd(ctx, "tags:post:1", "go", "redis", "backend")
	rdb.SAdd(ctx, "tags:post:2", "go", "docker", "backend")

	members, _ := rdb.SMembers(ctx, "tags:post:1").Result()
	fmt.Println("tags post:1 →", members)

	exists, _ := rdb.SIsMember(ctx, "tags:post:1", "redis").Result()
	fmt.Println("has 'redis' →", exists)

	common, _ := rdb.SInter(ctx, "tags:post:1", "tags:post:2").Result()
	fmt.Println("common tags (SINTER) →", common)

	all, _ := rdb.SUnion(ctx, "tags:post:1", "tags:post:2").Result()
	fmt.Println("all tags (SUNION) →", all)

	diff, _ := rdb.SDiff(ctx, "tags:post:1", "tags:post:2").Result()
	fmt.Println("exclusive to post:1 (SDIFF) →", diff)
}

func demoSortedSets(rdb *goredis.Client) {
	fmt.Println("\n=== SORTED SETS (Leaderboard) ===")
	rdb.Del(ctx, "leaderboard")

	rdb.ZAdd(ctx, "leaderboard",
		goredis.Z{Score: 1500, Member: "Alice"},
		goredis.Z{Score: 2300, Member: "Bob"},
		goredis.Z{Score: 1800, Member: "Charlie"},
		goredis.Z{Score: 3100, Member: "Diana"},
	)
	rdb.ZIncrBy(ctx, "leaderboard", 500, "Alice")

	top3, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
	fmt.Println("Top 3:")
	for i, z := range top3 {
		fmt.Printf("  #%d %-10s %.0f pts\n", i+1, z.Member, z.Score)
	}

	rank, _ := rdb.ZRevRank(ctx, "leaderboard", "Alice").Result()
	fmt.Printf("Alice rank → #%d\n", rank+1)

	score, _ := rdb.ZScore(ctx, "leaderboard", "Bob").Result()
	fmt.Printf("Bob score  → %.0f\n", score)
}

func demoHashes(rdb *goredis.Client) {
	fmt.Println("\n=== HASHES ===")
	rdb.Del(ctx, "user:100")

	rdb.HSet(ctx, "user:100",
		"name", "Alice",
		"email", "alice@example.com",
		"age", 30,
		"role", "admin",
	)

	name, _ := rdb.HGet(ctx, "user:100", "name").Result()
	fmt.Println("HGET name →", name)

	user, _ := rdb.HGetAll(ctx, "user:100").Result()
	fmt.Println("HGETALL →", user)

	rdb.HIncrBy(ctx, "user:100", "age", 1)
	age, _ := rdb.HGet(ctx, "user:100", "age").Int()
	fmt.Println("age after HINCRBY →", age)

	rdb.HDel(ctx, "user:100", "role")
	has, _ := rdb.HExists(ctx, "user:100", "role").Result()
	fmt.Println("still has role →", has)
}
