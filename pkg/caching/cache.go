package caching

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



type RedisClient struct{

	Db *redis.Client

}

func NewRedisClient(hostname string,password string,database int, protocol int) *RedisClient{

	rdb:= redis.NewClient(&redis.Options{
		Addr: hostname,
		Password: password,
		DB: database,
		Protocol: protocol,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return &RedisClient{
		Db: rdb,
	}
}





func (client *RedisClient) SetItem(ctx context.Context, key string, value interface{}, ttl time.Duration) error{
	return client.Db.Set(ctx,key,value,ttl).Err()
}

func(client *RedisClient) GetItem(ctx context.Context,key string) (string,error){
	return client.Db.Get(ctx,key).Result()
}

func (client *RedisClient) DeleteItem(ctx context.Context,key string) error{
	return client.Db.Del(ctx,key).Err()
}

func (client *RedisClient) ClearAll(ctx context.Context) error {
	return client.Db.FlushAll(ctx).Err()
}

func (client *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := client.Db.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := client.Db.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}

func (client *RedisClient) CloseConnection() error {
	return client.Db.Close()
}
