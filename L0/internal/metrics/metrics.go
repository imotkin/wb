package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Общее число HTTP-запросов",
	})

	ordersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_total",
		Help: "Общее число добавленных заказов",
	})

	failedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_total",
		Help: "Общее число ошибок при добавлении заказов",
	})
)

func IncRequests() {
	requestsTotal.Inc()
}

func IncOrders() {
	ordersTotal.Inc()
}

func IncFailed() {
	failedTotal.Inc()
}

func Handler() http.Handler {
	return promhttp.Handler()
}
