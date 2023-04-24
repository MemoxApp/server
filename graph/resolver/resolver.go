package resolver

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"time_speak_server/src/config"
	"time_speak_server/src/service/mail"
	"time_speak_server/src/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	conf    config.Config
	db      *mongo.Database
	r       *redis.Client
	userSvc *user.Svc
	mailSvc *mail.Svc
}

func NewResolver(conf config.Config, db *mongo.Database, r *redis.Client) *Resolver {
	return &Resolver{
		conf:    conf,
		db:      db,
		r:       r,
		userSvc: user.NewUserSvc(conf.User, db, r),
		mailSvc: mail.NewMailSvc(conf.Mail, r),
	}
}
