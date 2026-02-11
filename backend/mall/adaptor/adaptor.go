package adaptor

import (
	"mall/config"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdaptor interface {
	GetConfig() *config.Config
	GetRedis() *redis.Client
	GetDB() *gorm.DB
}

type Adaptor struct {
	conf  *config.Config
	db    *gorm.DB
	redis *redis.Client
}

func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		conf:  conf,
		db:    db,
		redis: redis,
	}
}

func (r *Adaptor) GetConfig() *config.Config {
	return r.conf
}
func (r *Adaptor) GetRedis() *redis.Client {
	return r.redis
}
func (r *Adaptor) GetDB() *gorm.DB {
	return r.db
}
