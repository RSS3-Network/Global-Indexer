package accesslog

import "time"

type AccessLog struct {
	ClientIP  string    `json:"client_ip"`
	Host      string    `json:"host"`
	URI       string    `json:"uri"`
	Consumer  *string   `json:"consumer"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"@timestamp"`
	RouteID   string    `json:"route_id"`
}
