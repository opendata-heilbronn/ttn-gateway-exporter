package exporter

import (
	"context"
	"fmt"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/config"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/logging"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/ttnclient"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
	"time"
)

var log = logging.Logger("target")

type Target struct {
	config config.Target
	client *ttnclient.TTNClient
	descs  map[string]*prometheus.Desc
}

func NewTarget(config config.Target) (*Target, error) {
	client, err := ttnclient.NewTTNClient(config.BaseUrl, ttnclient.ApiKeyAuthenticator{ApiKey: config.APIKey})
	if err != nil {
		return nil, err
	}
	return &Target{
		config: config,
		client: client,
		descs: map[string]*prometheus.Desc{
			"last_scrape_result":        desc(config.GatewayID, metricName("last_scrape_result"), "1 if the scrape from the TTN API was successful", []string{}),
			"connected_at":              desc(config.GatewayID, metricName("connected_at"), "Time the Gateway connected", []string{}),
			"disconnected_at":           desc(config.GatewayID, metricName("disconnected_at"), "Time the Gateway disconnected", []string{}),
			"last_status_at":            desc(config.GatewayID, metricName("last_status_at"), "Time TTN last received a status from the Gateway", []string{}),
			"last_uplink_at":            desc(config.GatewayID, metricName("last_uplink_at"), "Time TTN last received an uplink from the Gateway", []string{}),
			"last_downlink_at":          desc(config.GatewayID, metricName("last_downlink_at"), "Time TTN last sent a downlink to the Gateway", []string{}),
			"downlink_count":            desc(config.GatewayID, metricName("downlink_count"), "Number of downlinks through this Gateway", []string{}),
			"uplink_count":              desc(config.GatewayID, metricName("uplink_count"), "Number of uplinks through this Gateway", []string{}),
			"rtt_min":                   desc(config.GatewayID, metricName("rtt_min"), "Minimum round-trip-time", []string{}),
			"rtt_max":                   desc(config.GatewayID, metricName("rtt_max"), "Maximum round-trip-time", []string{}),
			"rtt_median":                desc(config.GatewayID, metricName("rtt_median"), "Median round-trip-time", []string{}),
			"rtt_count":                 desc(config.GatewayID, metricName("rtt_count"), "Number of round-trips", []string{}),
			"time":                      desc(config.GatewayID, metricName("time"), "Gateway time", []string{}),
			"boot_time":                 desc(config.GatewayID, metricName("boot_time"), "Gateway boot time", []string{}),
			"version":                   desc(config.GatewayID, metricName("version"), "Constantly 1. Exports the version of a subsystem as label.", []string{"subsystem", "version"}),
			"ip":                        desc(config.GatewayID, metricName("ip"), "Constantly 1. Exports the IP of the Gateway as label", []string{"num", "ip"}),
			"protocol":                  desc(config.GatewayID, metricName("protocol"), "Constantly 1. Exports the used protocol by the Gateway as label", []string{"protocol"}),
			"status_metrics":            desc(config.GatewayID, metricName("status_metrics"), "Gateway status metrics", []string{"metric"}),
			"antenna_location":          desc(config.GatewayID, metricName("antenna_location"), "Constantly 1. Antenna Location", []string{"antenna", "lat", "lon", "accuracy", "altitude", "source"}),
			"antenna_location_lat":      desc(config.GatewayID, metricName("antenna_location_lat"), "Antenna Latitude", []string{"antenna"}),
			"antenna_location_lon":      desc(config.GatewayID, metricName("antenna_location_lon"), "Antenna Longitude", []string{"antenna"}),
			"antenna_location_alt":      desc(config.GatewayID, metricName("antenna_location_alt"), "Antenna Altitude", []string{"antenna"}),
			"antenna_location_accuracy": desc(config.GatewayID, metricName("antenna_location_accuracy"), "Antenna location accuracy", []string{"antenna"}),
			"antenna_location_source":   desc(config.GatewayID, metricName("antenna_location_source"), "Constantly 1. Exports the antenna location source as label.", []string{"antenna", "source"}),
			"subband_utilization_limit": desc(config.GatewayID, metricName("subband_utilization_limit"), "Sub-band utilization limit", []string{"freqMin", "freqMax"}),
			"subband_utilization":       desc(config.GatewayID, metricName("subband_utilization"), "Sub-band utilization", []string{"freqMin", "freqMax"}),
		},
	}, nil
}

func (t *Target) Describe(descs chan<- *prometheus.Desc) {
	for _, desc := range t.descs {
		descs <- desc
	}
}

func (t *Target) Collect(metrics chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := t.client.GetGatewayConnectionStats(ctx, t.config.GatewayID)
	if err != nil {
		metrics <- prometheus.MustNewConstMetric(t.descs["last_scrape_result"], prometheus.GaugeValue, 0)
		log.Errorw("scrape error", "target", t.config.GatewayID, "error", err)
		return
	} else {
		metrics <- prometheus.MustNewConstMetric(t.descs["last_scrape_result"], prometheus.GaugeValue, 1)
	}

	if stats.DownlinkCount == "" {
		metrics <- prometheus.MustNewConstMetric(t.descs["downlink_count"], prometheus.CounterValue, 0)
	} else if downlinkCount, err := strconv.ParseInt(stats.DownlinkCount, 10, 64); err != nil {
		log.Errorw("numeric string to int conversion error", "target", t.config.GatewayID, "source", "downlink_count", "value", stats.DownlinkCount)
	} else {
		metrics <- prometheus.MustNewConstMetric(t.descs["downlink_count"], prometheus.CounterValue, float64(downlinkCount))
	}

	if stats.UplinkCount == "" {
		metrics <- prometheus.MustNewConstMetric(t.descs["uplink_count"], prometheus.CounterValue, 0)
	} else if uplinkCount, err := strconv.ParseInt(stats.UplinkCount, 10, 64); err != nil {
		log.Errorw("numeric string to int conversion error", "target", t.config.GatewayID, "source", "uplink_count", "value", stats.UplinkCount)
	} else {
		metrics <- prometheus.MustNewConstMetric(t.descs["uplink_count"], prometheus.CounterValue, float64(uplinkCount))
	}

	metrics <- prometheus.MustNewConstMetric(t.descs["connected_at"], prometheus.GaugeValue, unixTime(stats.ConnectedAt))
	metrics <- prometheus.MustNewConstMetric(t.descs["disconnected_at"], prometheus.GaugeValue, unixTime(stats.DisconnectedAt))
	metrics <- prometheus.MustNewConstMetric(t.descs["last_status_at"], prometheus.GaugeValue, unixTime(stats.LastStatusReceivedAt))
	metrics <- prometheus.MustNewConstMetric(t.descs["last_uplink_at"], prometheus.GaugeValue, unixTime(stats.LastUplinkReceivedAt))
	metrics <- prometheus.MustNewConstMetric(t.descs["last_downlink_at"], prometheus.GaugeValue, unixTime(stats.LastDownlinkReceivedAt))
	metrics <- prometheus.MustNewConstMetric(t.descs["rtt_min"], prometheus.GaugeValue, float64(stats.RoundTripTimes.Min))
	metrics <- prometheus.MustNewConstMetric(t.descs["rtt_max"], prometheus.GaugeValue, float64(stats.RoundTripTimes.Max))
	metrics <- prometheus.MustNewConstMetric(t.descs["rtt_median"], prometheus.GaugeValue, float64(stats.RoundTripTimes.Median))
	metrics <- prometheus.MustNewConstMetric(t.descs["rtt_count"], prometheus.CounterValue, float64(stats.RoundTripTimes.Count))
	metrics <- prometheus.MustNewConstMetric(t.descs["time"], prometheus.GaugeValue, unixTime(stats.LastStatus.Time))
	metrics <- prometheus.MustNewConstMetric(t.descs["boot_time"], prometheus.GaugeValue, unixTime(stats.LastStatus.BootTime))
	for subsystem, version := range stats.LastStatus.Versions {
		metrics <- prometheus.MustNewConstMetric(t.descs["version"], prometheus.GaugeValue, 1, subsystem, version)
	}
	for i, ip := range stats.LastStatus.IP {
		metrics <- prometheus.MustNewConstMetric(t.descs["ip"], prometheus.GaugeValue, 1, fmt.Sprintf("%d", i), ip)
	}
	metrics <- prometheus.MustNewConstMetric(t.descs["protocol"], prometheus.GaugeValue, 1, stats.Protocol)
	for metricName, metricValue := range stats.LastStatus.Metrics {
		metrics <- prometheus.MustNewConstMetric(t.descs["status_metrics"], prometheus.GaugeValue, float64(metricValue), metricName)
	}
	for i, antennaLocation := range stats.LastStatus.AntennaLocations {
		antennaNumber := fmt.Sprintf("%d", i)
		metrics <- prometheus.MustNewConstMetric(t.descs["antenna_location_lat"], prometheus.GaugeValue, antennaLocation.Latitude, antennaNumber)
		metrics <- prometheus.MustNewConstMetric(t.descs["antenna_location_lon"], prometheus.GaugeValue, antennaLocation.Longitude, antennaNumber)
		metrics <- prometheus.MustNewConstMetric(t.descs["antenna_location_alt"], prometheus.GaugeValue, float64(antennaLocation.Altitude), antennaNumber)
		metrics <- prometheus.MustNewConstMetric(t.descs["antenna_location_accuracy"], prometheus.GaugeValue, float64(antennaLocation.Accuracy), antennaNumber)
		metrics <- prometheus.MustNewConstMetric(t.descs["antenna_location_source"], prometheus.GaugeValue, 1, antennaNumber, antennaLocation.Source)
		metrics <- prometheus.MustNewConstMetric(
			t.descs["antenna_location"],
			prometheus.GaugeValue,
			1,
			antennaNumber,
			fmt.Sprintf("%f", antennaLocation.Latitude),
			fmt.Sprintf("%f", antennaLocation.Longitude),
			fmt.Sprintf("%d", antennaLocation.Altitude),
			fmt.Sprintf("%d", antennaLocation.Accuracy),
			antennaLocation.Source,
		)
	}
	for _, band := range stats.SubBands {
		metrics <- prometheus.MustNewConstMetric(t.descs["subband_utilization_limit"], prometheus.GaugeValue, band.DownlinkUtilizationLimit, band.MinFrequency, band.MaxFrequency)
		metrics <- prometheus.MustNewConstMetric(t.descs["subband_utilization"], prometheus.GaugeValue, band.DownlinkUtilization, band.MinFrequency, band.MaxFrequency)
	}
}

func metricName(names ...string) string {
	return prometheus.BuildFQName("ttn", "gateway", strings.Join(names, "_"))
}

func desc(gw, name, help string, variableLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(name, help, variableLabels, prometheus.Labels{
		"gateway": gw,
	})
}

func unixTime(in time.Time) float64 {
	if in.IsZero() {
		return 0
	}
	return float64(in.Unix())
}
