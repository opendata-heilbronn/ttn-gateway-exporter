package ttnclient

import (
	"time"
)

// GatewayConnectionStats https://www.thethingsindustries.com/docs/reference/api/gateway_server/#message:GatewayConnectionStats
type GatewayConnectionStats struct {
	ConnectedAt            time.Time      `json:"connected_at"`
	DisconnectedAt         time.Time      `json:"disconnected_at"`
	Protocol               string         `json:"protocol"`
	LastStatusReceivedAt   time.Time      `json:"last_status_received_at"`
	LastStatus             GatewayStatus  `json:"last_status"`
	LastUplinkReceivedAt   time.Time      `json:"last_uplink_received_at"`
	UplinkCount            string         `json:"uplink_count"`
	LastDownlinkReceivedAt time.Time      `json:"last_downlink_received_at"`
	DownlinkCount          string         `json:"downlink_count"`
	RoundTripTimes         RoundTripTimes `json:"round_trip_times"`
	SubBands               []SubBand      `json:"sub_bands"`
}

// SubBand https://www.thethingsindustries.com/docs/reference/api/gateway_server/#message:GatewayConnectionStats.RoundTripTimes
type SubBand struct {
	MinFrequency             string  `json:"min_frequency"`
	MaxFrequency             string  `json:"max_frequency"`
	DownlinkUtilizationLimit float64 `json:"downlink_utilization_limit"`
	DownlinkUtilization      float64 `json:"downlink_utilization,omitempty"`
}

// RoundTripTimes https://www.thethingsindustries.com/docs/reference/api/gateway_server/#message:GatewayConnectionStats.RoundTripTimes
type RoundTripTimes struct {
	Min    time.Duration `json:"min"`
	Max    time.Duration `json:"max"`
	Median time.Duration `json:"median"`
	Count  uint32        `json:"count"`
}

// GatewayStatus https://www.thethingsindustries.com/docs/reference/api/gateway_server/#message:GatewayStatus
type GatewayStatus struct {
	Time             time.Time              `json:"time"`
	BootTime         time.Time              `json:"boot_time"`
	Versions         map[string]string      `json:"versions"`
	AntennaLocations []Location             `json:"antenna_locations"`
	IP               []string               `json:"ip"`
	Metrics          map[string]float64     `json:"metrics"`
	Advanced         map[string]interface{} `json:"advanced"`
}

// Location https://www.thethingsindustries.com/docs/reference/api/gateway_server/#message:Location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  int32   `json:"altitude"`
	Accuracy  int32   `json:"accuracy"`

	// https://www.thethingsindustries.com/docs/reference/api/gateway_server/#enum:LocationSource
	Source string `json:"source"`
}
