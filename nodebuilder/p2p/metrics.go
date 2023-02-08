package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/metrics"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
)

// global meter provider (see opentelemetry docs)
var (
	meter = global.MeterProvider().Meter("p2p")
)

// WithMetrics option sets up metrics for p2p networking.
func WithMetrics(bc *metrics.BandwidthCounter) {
	bandwidthTotalInbound, _ := meter.
		SyncInt64().
		Histogram(
			"p2p_bandwidth_total_inbound",
			instrument.WithUnit(unit.Bytes),
			instrument.WithDescription("total number of bytes received by the host"),
		)
	bandwidthTotalOutbound, _ := meter.
		SyncInt64().
		Histogram(
			"p2p_bandwidth_total_outbound",
			instrument.WithUnit(unit.Bytes),
			instrument.WithDescription("total number of bytes sent by the host"),
		)
	bandwidthRateInbound, _ := meter.
		SyncFloat64().
		Histogram(
			"p2p_bandwidth_rate_inbound",
			instrument.WithDescription("total number of bytes sent by the host"),
		)
	bandwidthRateOutbound, _ := meter.
		SyncFloat64().
		Histogram(
			"p2p_bandwidth_rate_outbound",
			instrument.WithDescription("total number of bytes sent by the host"),
		)
	p2pPeerCount, _ := meter.
		AsyncFloat64().
		Gauge(
			"p2p_peer_count",
			instrument.WithDescription("number of peers connected to the host"),
		)

	err := meter.RegisterCallback(
		[]instrument.Asynchronous{
			p2pPeerCount,
		}, func(ctx context.Context) {
			bcStats := bc.GetBandwidthTotals()
			bcByPeerStats := bc.GetBandwidthByPeer()

			bandwidthTotalInbound.Record(ctx, bcStats.TotalIn)
			bandwidthTotalOutbound.Record(ctx, bcStats.TotalOut)
			bandwidthRateInbound.Record(ctx, bcStats.RateIn)
			bandwidthRateOutbound.Record(ctx, bcStats.RateOut)

			p2pPeerCount.Observe(ctx, float64(len(bcByPeerStats)))
		},
	)

	if err != nil {
		panic(err)
	}
}
