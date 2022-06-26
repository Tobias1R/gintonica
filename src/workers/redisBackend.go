package workers

import (
	"github.com/go-redis/redis/v8"
)

func GetRedisClient() redis.Client {
	return getRedisClient()
}

func getRedisClient() redis.Client {
	//ctx := context.Background()

	//redis://redis:eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81@172.20.0.5:6379
	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.20.0.5:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // no password set
		DB:       0,                                  // use default DB
	})

	return *rdb

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)
}
