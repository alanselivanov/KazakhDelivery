package application

import (
	"log"
	"sync/atomic"
)

type Metrics struct {
	eventsProcessedTotal   uint64
	stockUpdateErrorsTotal uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncEventsProcessed() {
	atomic.AddUint64(&m.eventsProcessedTotal, 1)
	log.Printf("Metrics: events_processed_total=%d", m.eventsProcessedTotal)
}

func (m *Metrics) IncStockUpdateErrors() {
	atomic.AddUint64(&m.stockUpdateErrorsTotal, 1)
	log.Printf("Metrics: stock_update_errors_total=%d", m.stockUpdateErrorsTotal)
}

func (m *Metrics) GetEventsProcessedTotal() uint64 {
	return atomic.LoadUint64(&m.eventsProcessedTotal)
}

func (m *Metrics) GetStockUpdateErrorsTotal() uint64 {
	return atomic.LoadUint64(&m.stockUpdateErrorsTotal)
}
