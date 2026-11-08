package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "restreamer"

const (
	labelStream = "stream"
	labelTarget = "target"
	labelKind   = "kind"
	labelCodec  = "codec"
	labelCode   = "code"
)

var registry = prometheus.NewRegistry()

var factory = promauto.With(registry)

func init() {
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}

// Handler serves all registered metrics in the Prometheus exposition format.
func Handler() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})
}

var (
	RTMPSessions = factory.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "rtmp_sessions",
		Help:      "Number of open incoming RTMP connections (publishers and players).",
	})

	IngestActive = factory.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "ingest_active",
		Help:      "1 if a publisher is currently sending the stream, 0 otherwise.",
	}, []string{labelStream})

	IngestFrames = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "ingest_frames_total",
		Help:      "Frames received from the publisher.",
	}, []string{labelStream, labelKind, labelCodec})

	IngestBytes = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "ingest_bytes_total",
		Help:      "Payload bytes received from the publisher.",
	}, []string{labelStream, labelKind, labelCodec})

	IngestKeyframes = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "ingest_keyframes_total",
		Help:      "Video keyframes (IDR) received from the publisher.",
	}, []string{labelStream})

	StreamDroppedBatches = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "stream_dropped_batches_total",
		Help:      "Frame batches dropped because the stream dispatch buffer was full.",
	}, []string{labelStream})

	PushConnected = factory.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "push_connected",
		Help:      "1 if the push target accepted the publish, 0 otherwise.",
	}, []string{labelStream, labelTarget})

	PushConnectFailures = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "push_connect_failures_total",
		Help:      "Push target sessions that ended with an error (each one is followed by a reconnect).",
	}, []string{labelStream, labelTarget})

	PushErrors = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "push_errors_total",
		Help:      "RTMP error statuses reported by the push target, by status code.",
	}, []string{labelStream, labelTarget, labelCode})

	PushBytes = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "push_bytes_total",
		Help:      "Frame payload bytes sent to the push target.",
	}, []string{labelStream, labelTarget})

	PushQueueBatches = factory.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "push_queue_batches",
		Help:      "Frame batches waiting to be sent to the push target.",
	}, []string{labelStream, labelTarget})

	PushDroppedBatches = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "push_dropped_batches_total",
		Help:      "Frame batches dropped because the queue of a connected push target overflowed.",
	}, []string{labelStream, labelTarget})

	PullConsumers = factory.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "pull_consumers",
		Help:      "Players currently watching the stream directly from the restreamer.",
	}, []string{labelStream})

	PullDroppedBatches = factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "pull_dropped_batches_total",
		Help:      "Frame batches dropped because a player queue overflowed.",
	}, []string{labelStream})
)

var pushVecs = []*prometheus.MetricVec{
	PushConnected.MetricVec,
	PushConnectFailures.MetricVec,
	PushErrors.MetricVec,
	PushBytes.MetricVec,
	PushQueueBatches.MetricVec,
	PushDroppedBatches.MetricVec,
}

var streamVecs = append([]*prometheus.MetricVec{
	IngestActive.MetricVec,
	IngestFrames.MetricVec,
	IngestBytes.MetricVec,
	IngestKeyframes.MetricVec,
	StreamDroppedBatches.MetricVec,
	PullConsumers.MetricVec,
	PullDroppedBatches.MetricVec,
}, pushVecs...)

// PushLabels returns the label set identifying a push target series.
func PushLabels(stream, target string) prometheus.Labels {
	return prometheus.Labels{labelStream: stream, labelTarget: target}
}

// DeletePushTarget removes all series of a push target, so that a removed
// target does not stay in the exported metrics forever.
func DeletePushTarget(stream, target string) {
	labels := PushLabels(stream, target)
	for _, v := range pushVecs {
		v.DeletePartialMatch(labels)
	}
}

// DeleteStream removes all series of a deleted stream.
func DeleteStream(stream string) {
	labels := prometheus.Labels{labelStream: stream}
	for _, v := range streamVecs {
		v.DeletePartialMatch(labels)
	}
}
