package bce

import "time"

// CallbackRequest 百度云回调
type CallbackRequest struct {
	Events []CallbackEvents `json:"events"`
}

type CallbackEvents struct {
	Version     string          `json:"version"`
	EventId     string          `json:"eventId"`
	EventOrigin string          `json:"eventOrigin"`
	EventTime   time.Time       `json:"eventTime"`
	EventType   string          `json:"eventType"`
	Content     CallbackContent `json:"content"`
}

type CallbackContent struct {
	UserId       string    `json:"userId"`
	OwnerId      string    `json:"ownerId"`
	AccessKeyId  string    `json:"accessKeyId"`
	Domain       string    `json:"domain"`
	Bucket       string    `json:"bucket"`
	Object       string    `json:"object"`
	Etag         string    `json:"etag"`
	ContentType  string    `json:"contentType"`
	Filesize     int       `json:"filesize"`
	LastModified time.Time `json:"lastModified"`
	Credentials  struct {
		AccessKeyId     string    `json:"accessKeyId"`
		SecretAccessKey string    `json:"secretAccessKey"`
		SessionToken    string    `json:"sessionToken"`
		Expiration      time.Time `json:"expiration"`
	} `json:"credentials"`
}
