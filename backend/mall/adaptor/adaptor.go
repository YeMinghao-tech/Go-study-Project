package adaptor

import (
	"mall/config"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Adaptor struct {
	Conf     *config.Config
	DbClient *gorm.DB
	Redis    *redis.Client
}

func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		Conf:     conf,
		DbClient: db,
		Redis:    redis,
	}
}
