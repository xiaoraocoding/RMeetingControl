package initialize

// 初始化redis
import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"sync"
	"time"
)

// RedisClient Redis 服务
type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// once 确保全局的 Redis 对象只实例一次
var once sync.Once

// Redis 全局 Redis，使用 db 1
var Redis *RedisClient

// ConnectRedis 连接 redis 数据库，设置全局的 Redis 对象
func ConnectRedis(address string, password string, db int) {
	once.Do(func() {
		Redis = NewClient(address, password, db)
	})
}

// NewClient 创建一个新的 redis 连接
func NewClient(address string, password string, db int) *RedisClient {

	// 初始化自定的 RedisClient 实例
	rds := &RedisClient{}
	// 使用默认的 context
	rds.Context = context.Background()

	// 使用 redis 库里的 NewClient 初始化连接
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// 测试一下连接
	err := rds.Ping()
	Log.Error("Error connect redis failed", err)
	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
		Log.Error("Error redis set failed", err)
		return false
	}
	return true
}

// Get 获取 key 对应的 value
func (rds RedisClient) Get(key string) string {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			Log.Error("Error redis get failed", err)
		}
		return ""
	}
	return result
}

func (rds RedisClient) Hset(key string, field string, value string) bool {
	if err := rds.Client.HSet(rds.Context, key, field, value).Err(); err != nil {
		Log.Error("Error redis Hset failed", err)
		return false
	}
	return true
}

func (rds RedisClient) Hget(key string, field string) string {
	result, err := rds.Client.HGet(rds.Context, key, field).Result()
	if err != nil {
		if err != redis.Nil {
			Log.Error("Error redis Hget failed", err)
		}
		return ""
	}
	return result
}

func (rds RedisClient) HgetAll(key string) map[string]string {
	result, err := rds.Client.HGetAll(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			Log.Error("Error redis HgetAll failed", err)
		}
		return nil
	}
	return result
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds RedisClient) Has(key string) bool {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			Log.Error("Error redis has failed", err)
		}
		return false
	}
	return true
}

func (rds RedisClient) HDel(key string, fields ...string) bool {
	if err := rds.Client.HDel(rds.Context, key, fields...).Err(); err != nil {
		Log.Error("Error redis Hdel failed", err)
		return false
	}
	return true
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds RedisClient) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Context, keys...).Err(); err != nil {
		Log.Error("Error redis delete failed", err)
		return false
	}
	return true
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds RedisClient) FlushDB() bool {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		Log.Error("Error redis flushDb failed", err)
		return false
	}
	return true
}

// Increment 当参数只有 1 个时，为 key，其值增加 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
func (rds RedisClient) Increment(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Incr(rds.Context, key).Err(); err != nil {
			Log.Error("Error redis Increment failed", err)
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.IncrBy(rds.Context, key, value).Err(); err != nil {
			Log.Error("Error redis Increment failed", err)
			return false
		}
	default:
		Log.Error("Error redis Increment 参数过多")
		return false
	}
	return true
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
func (rds RedisClient) Decrement(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Decr(rds.Context, key).Err(); err != nil {
			Log.Error("Error redis Decrement failed", err)
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.DecrBy(rds.Context, key, value).Err(); err != nil {
			Log.Error("Error redis Decrement failed", err)
			return false
		}
	default:
		Log.Error("Error redis Decrement 参数过多")
		return false
	}
	return true
}
