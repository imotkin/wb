package broker

import (
	"net"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Config struct {
	Host        string        `koanf:"host"`
	Port        string        `koanf:"port"`
	Topic       string        `koanf:"topic"`
	TopicDLQ    string        `koanf:"topic_dlq"`
	GroupID     string        `koanf:"group_id"`
	FirstOffset bool          `koanf:"first_offset"`
	Interval    time.Duration `koanf:"interval"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.Topic, validation.Required),
		validation.Field(&c.TopicDLQ, validation.Required),
		validation.Field(&c.GroupID, validation.Required),
		validation.Field(&c.FirstOffset, validation.In(true, false)),
		validation.Field(&c.Interval, validation.Required),
	)
}

func (c *Config) Endpoint() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *Config) Offset() kgo.Offset {
	if c.FirstOffset {
		return kgo.NewOffset().AtStart()
	}

	return kgo.NewOffset().AtEnd()
}
