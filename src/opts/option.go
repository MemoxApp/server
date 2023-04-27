package opts

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Option func(bson.M) bson.M

// With 设置键值对
func With(key, t string) Option {
	return func(m bson.M) bson.M {
		m[key] = t
		return m
	}
}

// WithTitle 设置标题
func WithTitle(t string) Option {
	return func(m bson.M) bson.M {
		m["title"] = t
		return m
	}
}

// WithArchived 设置归档
func WithArchived(t bool) Option {
	return func(m bson.M) bson.M {
		m["archived"] = t
		return m
	}
}

// WithEnabled 设置启用
func WithEnabled(t bool) Option {
	return func(m bson.M) bson.M {
		m["enabled"] = t
		return m
	}
}

// WithContent 设置内容
func WithContent(t string) Option {
	return func(m bson.M) bson.M {
		m["content"] = t
		return m
	}
}

// WithTags 设置话题
func WithTags(t []primitive.ObjectID) Option {
	return func(m bson.M) bson.M {
		m["hash_tags"] = t
		return m
	}
}
