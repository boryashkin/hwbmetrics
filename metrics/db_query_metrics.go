package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const MetricsSubsystemDbQuery = "db_query"

type IDbQueryMetrics interface {
	QueryDuration(time.Duration)
	IncDone()
	IncError(error)
}

func NewDbQueryMetrics(dbName, queryName string) IDbQueryMetrics {
	labels := map[string]string{
		MetricsLabelDatabase:      dbName,
		MetricsLabelDatabaseQuery: queryName,
	}

	return &dbRequestMetrics{
		durationMs: NewHistogram(MetricsNamespace, MetricsSubsystemDbQuery, "duration_ms", labels, DefaultDurationMsBuckets),
		nbDone:     NewCounter(MetricsNamespace, MetricsSubsystemDbQuery, "nb_done", labels),
		nbError:    NewCounter(MetricsNamespace, MetricsSubsystemDbQuery, "nb_error", labels),
	}
}

type dbRequestMetrics struct {
	durationMs prometheus.Histogram
	nbDone     prometheus.Counter
	nbError    prometheus.Counter
}

func (m *dbRequestMetrics) QueryDuration(duration time.Duration) {
	m.durationMs.Observe(float64(duration.Milliseconds()))
}
func (m *dbRequestMetrics) IncDone() {
	m.nbDone.Inc()
}
func (m *dbRequestMetrics) IncError(error) {
	m.nbError.Inc()
}

/*
Helper для быстрого расчета метрик запроса. Плюсом идет то, что метод сам анализирует ошибку и может инкрементить
нужные вспомогательные метрики.

Пример:

```

	func SomeDatabaseMethod() (e error){
		defer func(from time.Time) {
			metrics.DbQueryMetricsHelper(metrics, from, e)
		}(time.Now())

		// Your code goes here
	}

```
*/
func DbQueryMetricsHelper(m IDbQueryMetrics, startTm time.Time, err error) {
	m.QueryDuration(time.Since(startTm))
	if err != nil {
		m.IncError(err)
	} else {
		m.IncDone()
	}
}
