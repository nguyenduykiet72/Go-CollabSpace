package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ActiveConnections tracks the number of active WebSocket connections per document
	ActiveConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ws_active_connections",
		Help: "The total number of active WebSocket connections",
	}, []string{"doc_id"})

	// ProcessedMessages tracks the number of processed WebSocket messages by type
	ProcessedMessages = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ws_processed_messages_total",
		Help: "The total number of processed WebSocket messages",
	}, []string{"type"})
)
