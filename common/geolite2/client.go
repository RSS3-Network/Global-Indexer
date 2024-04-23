package geolite2

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/maxmind/geoipupdate/v6/pkg/geoipupdate"
	"github.com/oschwald/geoip2-golang"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
	"go.uber.org/zap"
)

type Client struct {
	reader *geoip2.Reader
}

func (c *Client) LookupNodeLocation(_ context.Context, endpoint string) ([]*schema.NodeLocation, error) {
	if c == nil || c.reader == nil {
		zap.L().Warn("geoip2 client is nil")

		return nil, nil
	}

	ips := make([]net.IP, 0)

	zap.L().Info("Looking up Node location", zap.String("endpoint", endpoint))

	if ip := net.ParseIP(endpoint); ip == nil {
		ipAddresses, err := net.LookupIP(endpoint)
		if err != nil {
			return nil, fmt.Errorf("lookup endpoint: %w", err)
		}

		for _, ipAddress := range ipAddresses {
			if ipv4 := ipAddress.To4(); ipv4 != nil {
				ips = append(ips, ipv4)
			}
		}
	} else {
		ips = append(ips, ip)
	}

	records := make([]*schema.NodeLocation, 0, len(ips))

	for _, ip := range ips {
		record, err := c.reader.City(ip)
		if err != nil {
			return nil, fmt.Errorf("get city: %w", err)
		}

		if record.Location.Longitude == 0 && record.Location.Latitude == 0 {
			continue
		}

		local := &schema.NodeLocation{
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
		}

		if len(record.Country.Names) > 0 {
			local.Country = record.Country.Names["en"]
		}

		if len(record.Subdivisions) > 0 && len(record.Subdivisions[0].Names) > 0 {
			local.Region = record.Subdivisions[0].Names["en"]
		}

		if len(record.City.Names) > 0 {
			local.City = record.City.Names["en"]
		}

		records = append(records, local)
	}

	return records, nil
}

func NewClient(conf *config.GeoIP) *Client {
	dir := filepath.Dir(conf.File)

	config := &geoipupdate.Config{
		URL:               "https://updates.maxmind.com",
		DatabaseDirectory: dir,
		LockFile:          filepath.Join(dir, ".geoipupdate.lock"),
		AccountID:         conf.Account,
		LicenseKey:        conf.LicenseKey,
		EditionIDs:        []string{"GeoLite2-City"},
		Output:            true,
		Verbose:           true,
		Parallelism:       1,
	}

	client := geoipupdate.NewClient(config)

	err := client.Run(context.Background())
	if err != nil {
		zap.L().Warn("run geoipupdate failed", zap.Error(err))
	}

	reader, err := geoip2.Open(conf.File)
	if err != nil {
		zap.L().Warn("open geoip2 database failed", zap.Error(err))
		return nil
	}

	return &Client{
		reader: reader,
	}
}
