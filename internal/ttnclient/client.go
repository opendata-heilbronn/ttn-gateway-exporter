package ttnclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

var log = logging.Logger("ttn-client")

var ttnApiRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "ttnapi",
	Subsystem: "client",
	Name:      "request_duration_seconds",
	Help:      "Histogram of the request duration towards the TTN API",
	Buckets:   []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10},
}, []string{"code", "method"})

var ttnApiRequestsInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "ttnapi",
	Subsystem: "client",
	Name:      "request_inflight",
	Help:      "Number of requests towards the TTN API that are currently ongoin",
})

var ttnRateLimitAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "ttnapi",
	Subsystem: "client",
	Name:      "ratelimit_allowed",
	Help:      "The maximum number of requests allowed by the TTN rate limiting",
})

var ttnRateLimitCurrent = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "ttnapi",
	Subsystem: "client",
	Name:      "ratelimit_current",
	Help:      "The maximum number of requests allowed by the TTN rate limiting",
})

func init() {
	prometheus.MustRegister(
		ttnApiRequestDuration,
		ttnApiRequestsInFlight,
		ttnRateLimitAllowed,
		ttnRateLimitCurrent,
	)
}

type TTNClient struct {
	baseUrl       url.URL
	authenticator Authenticator
	http          http.Client
}

func NewTTNClient(baseUrl string, authenticator Authenticator) (*TTNClient, error) {
	parsedUrl, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}

	return &TTNClient{
		baseUrl:       *parsedUrl,
		authenticator: authenticator,
		http: http.Client{
			Timeout: 10 * time.Second,
			Transport: promhttp.InstrumentRoundTripperDuration(
				ttnApiRequestDuration,
				promhttp.InstrumentRoundTripperInFlight(
					ttnApiRequestsInFlight,
					http.DefaultTransport,
				),
			),
		},
	}, nil
}

func (client *TTNClient) GetGatewayConnectionStats(ctx context.Context, gatewayId string) (stats GatewayConnectionStats, err error) {
	reqUrl := client.baseUrl
	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/api/v3/gs/gateways/%s/connection/stats", gatewayId))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return stats, err
	}

	err = client.authenticator.Authenticate(req)
	if err != nil {
		return stats, err
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return stats, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if warnHeader := resp.Header.Get("x-warning"); warnHeader != "" {
		log.Warnw("ttn api warning", "content", warnHeader)
	}
	if rateLimitAvailable := resp.Header.Get("x-rate-limit-available"); rateLimitAvailable != "" {
		rateLimitAvailableNum, err := strconv.ParseInt(rateLimitAvailable, 10, 64)
		if err == nil {
			ttnRateLimitCurrent.Set(float64(rateLimitAvailableNum))
		}
	}
	if rateLimitAllowed := resp.Header.Get("x-rate-limit-limit"); rateLimitAllowed != "" {
		rateLimitAllowedNum, err := strconv.ParseInt(rateLimitAllowed, 10, 64)
		if err == nil {
			ttnRateLimitAllowed.Set(float64(rateLimitAllowedNum))
		}
	}

	if resp.StatusCode != 200 {
		respBuf, err := io.ReadAll(resp.Body)
		if err != nil {
			return stats, err
		}
		var potentialJsonBody map[string]interface{}
		err = json.Unmarshal(respBuf, &potentialJsonBody)
		if err != nil {
			return stats, fmt.Errorf("TTN API responded with non 200 status code %s: %s", resp.Status, string(respBuf))
		}
		return stats, fmt.Errorf("TTN API responded with non 200 status code %s: %#v", resp.Status, potentialJsonBody)
	}

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
