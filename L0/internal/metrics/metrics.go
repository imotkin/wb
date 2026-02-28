package metrics

import (
	"net/http"

	"github.com/imotkin/L0/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	counters map[string]prometheus.Counter
	gauges   map[string]prometheus.Gauge
}

func New(log logger.Logger) (Metrics, error) {
	counters := map[string]prometheus.Counter{
		"RequestsTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Общее число HTTP-запросов",
		}),

		"OrdersTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Общее число добавленных заказов",
		}),

		"FailedTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "failed_total",
			Help: "Общее число ошибок при добавлении заказов",
		}),

		"CacheGetTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_get_total",
			Help: "Общее число полученных заказов из кэша",
		}),

		"CacheSetTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_set_total",
			Help: "Общее число добавленных заказов в кэш",
		}),

		"PostgresGetTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "pg_get_total",
			Help: "Общее число полученных заказов из базы данных",
		}),

		"PostgresSetTotal": promauto.NewCounter(prometheus.CounterOpts{
			Name: "pg_set_total",
			Help: "Общее число добавленных заказов в базу данных",
		}),
	}

	gauges := map[string]prometheus.Gauge{
		"KafkaStatus": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "kafka_status",
			Help: "Текущий статус для Kafka (1 - доступна, 0 - нет)",
		}),

		"PostgresStatus": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "postgres_status",
			Help: "Текущий статус для PostgreSQL (1 - доступна, 0 - нет)",
		}),
	}

	return &metrics{
		counters: counters,
		gauges:   gauges,
	}, nil
}

func (m *metrics) IncCounter(metric string) {
	m.counters[metric].Inc()
}

func (m *metrics) IncRequests() {
	m.IncCounter("RequestsTotal")
}

func (m *metrics) IncOrders() {
	m.IncCounter("OrdersTotal")
}

func (m *metrics) IncFailed() {
	m.IncCounter("FailedTotal")
}

func (m *metrics) IncCacheGet() {
	m.IncCounter("CacheGetTotal")
}

func (m *metrics) IncCacheSet() {
	m.IncCounter("CacheSetTotal")
}

func (m *metrics) IncPostgresGet() {
	m.IncCounter("PostgresGetTotal")
}

func (m *metrics) IncPostgresSet() {
	m.IncCounter("PostgresSetTotal")
}

func (m *metrics) SetKafkaStatus(i int) {
	m.gauges["KafkaStatus"].Set(float64(i))
}

func (m *metrics) SetPostgresStatus(i int) {
	m.gauges["PostgresStatus"].Set(float64(i))
}

func Handler() http.Handler {
	return promhttp.Handler()
}
