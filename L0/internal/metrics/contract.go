package metrics

type Metrics interface {
	IncRequests()
	IncOrders()
	IncFailed()
	IncCacheGet()
	IncCacheSet()
	IncPostgresGet()
	IncPostgresSet()
	SetKafkaStatus(int)
	SetPostgresStatus(int)
}
