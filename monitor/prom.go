package monitor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gomod.usaken.org/uvcounter/spine"
)

var EnableMetric = true

func RunPrometheusServer() error {
	spine.SystemGroup.Add(1)
	defer spine.SystemGroup.Done()

	if !EnableMetric {
		return nil
	}

	prometheus.Register(counter)
	prometheus.Register(responseLatencyHistogram)
	prometheus.Register(responseLatencySummary)

	middleware := httpmiddleware.New(prometheus.DefaultRegisterer, nil)
	metricHandler := middleware.WrapHandler(
		"/metrics",
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{},
		),
	)

	server := http.Server{Addr: ":9000", Handler: metricHandler}

	go func() {
		spine.SystemGroup.Add(1)
		defer spine.SystemGroup.Done()

		reason := <-spine.C.Done()
		fmt.Printf("prom server shutdown started due to %s\n", reason)
		// 5분은 휴리스틱하게 정해진 시간이다.
		// prometheus 서버를 내리기전에 이미 충분하게 요청이 들어오지 않은 상태이겠지만,
		// 혹시 1분이상 실행중인 요청이 있다면 실패하도록한다.
		// timeout값보다 크게 하여 최대한 보수적으로 잡는다.
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		err := server.Shutdown(ctx)
		if err != nil {
			fmt.Printf("prom server shutdown failed %e\n", err)
		}
		fmt.Printf("prom server successfully shutdown\n")
	}()

	go func() {
		spine.SystemGroup.Add(1)
		defer spine.SystemGroup.Done()

		server.ListenAndServe()
		fmt.Printf("prom server shutdown end\n")
	}()

	fmt.Printf("prom server running... \n")

	return nil
}

var counter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace:   "ic",
		Subsystem:   "http",
		Name:        "requests",
		Help:        "",
		ConstLabels: nil,
	},
	[]string{"uri"},
)

var responseLatencyHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "ic",
		Subsystem: "http",
		Name:      "response_latency",
		Buckets:   prometheus.LinearBuckets(0, 20, 2000),
	},
	[]string{"uri"},
)

var responseLatencySummary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "ic",
		Subsystem: "http",
		Name:      "response_latency_summary",
		Objectives: map[float64]float64{
			0.99999: 0.001,
			0.999:   0.001,
			0.99:    0.001,
			0.95:    0.001,
			0.90:    0.001,
			0.80:    0.001,
			0.70:    0.001,
			0.60:    0.001,
			0.50:    0.001,
			0.40:    0.001,
			0.30:    0.001,
			0.20:    0.001,
			0.10:    0.001,
		},
	},
	[]string{"uri"},
)

func CollectHTTPRequest(uri string) {
	counter.WithLabelValues(uri).Inc()
}

func CollectHTTPResponse(uri string, code int, d time.Duration) {
	duration := float64(d.Nanoseconds())
	responseLatencyHistogram.WithLabelValues(uri).Observe(duration)
	responseLatencySummary.WithLabelValues(uri).Observe(duration)
}
