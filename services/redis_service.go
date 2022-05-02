package services

import (
	"github.com/go-redis/redis"
	"time"
)

type IRedisService interface {
	Connect() *redis.Client
	Get(key string) []byte
	Set(key string, value interface{}, expiration time.Duration) error
}

type RedisService struct {
	Option redis.Options
	Client *redis.Client
}

func NewRedisService(options redis.Options) *RedisService {
	return &RedisService{
		Option: options,
	}
}

func (r *RedisService) Connect() *redis.Client {
	r.Client = redis.NewClient(&r.Option)

	return r.Client
}

func (r *RedisService) Get(key string) ([]byte, error) {
	value := r.Client.Get(key)
	bytes, err := value.Bytes()

	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	cmd := r.Client.Set(key, value, expiration)
	return cmd.Err()
}
