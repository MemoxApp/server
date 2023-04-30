package resolver

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"memox_server/src/config"
	"memox_server/src/service/comment"
	"memox_server/src/service/hashtag"
	"memox_server/src/service/history"
	"memox_server/src/service/mail"
	"memox_server/src/service/memory"
	"memox_server/src/service/resource"
	"memox_server/src/service/storage"
	"memox_server/src/service/storage/utils"
	"memox_server/src/service/subscribe"
	"memox_server/src/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	conf         config.Config
	db           *mongo.Database
	r            *redis.Client
	userSvc      *user.Svc
	mailSvc      *mail.Svc
	memorySvc    *memory.Svc
	hashtagSvc   *hashtag.Svc
	historySvc   *history.Svc
	commentSvc   *comment.Svc
	subscribeSvc *subscribe.Svc
	resourceSvc  *resource.Svc
	storageSvc   utils.Service
}

func NewResolver(conf config.Config, db *mongo.Database, r *redis.Client) *Resolver {
	sto := storage.NewStorageSvc(conf.Storage, r)
	return &Resolver{
		conf:         conf,
		db:           db,
		r:            r,
		userSvc:      user.NewUserSvc(conf.User, db, r),
		mailSvc:      mail.NewMailSvc(conf.Mail, r),
		memorySvc:    memory.NewMemorySvc(conf.Memory, db, r),
		hashtagSvc:   hashtag.NewHashTagSvc(conf.Hashtag, db, r),
		historySvc:   history.NewHistorySvc(conf.History, db, r),
		commentSvc:   comment.NewCommentSvc(conf.Comment, db, r),
		subscribeSvc: subscribe.NewSubscribeSvc(conf.Subscribe, db, r),
		storageSvc:   sto,
		resourceSvc:  resource.NewResourceSvc(conf.Resource, db, r, sto),
	}
}
