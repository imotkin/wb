package metrics

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type Endpoint struct {
	URL      string
	Name     string
	Help     string
	Interval time.Duration
	metric   prometheus.Gauge
}

func (e Endpoint) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.URL, validation.Required),
		validation.Field(&e.Name, validation.Required),
		validation.Field(&e.Help, validation.Required),
		validation.Field(&e.Interval, validation.Required),
	)
}
