package ingestion

import (
	"go.opentelemetry.io/otel"
)

var (
	meter = otel.Meter("services/internal/ingestion")

	locationUpdatesReceivedTotal, _ = meter.Int64Counter("location_updates_received_total")
	locationUpdatesInvalidTotal, _  = meter.Int64Counter("location_updates_invalid_total")
)
