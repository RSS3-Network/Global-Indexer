package geolite2

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/maxmind/geoipupdate/v6/pkg/geoipupdate"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/oschwald/geoip2-golang"
)

type Client struct {
	reader *geoip2.Reader
}

func (c *Client) LookupLocal(_ context.Context, endpoint string) ([]*schema.NodeLocal, error) {
	ips := make([]net.IP, 0)

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

	records := make([]*schema.NodeLocal, 0, len(ips))

	for _, ip := range ips {
		record, err := c.reader.City(ip)
		if err != nil {
			return nil, fmt.Errorf("get city: %w", err)
		}

		if record.Location.Longitude == 0 && record.Location.Latitude == 0 {
			continue
		}

		local := &schema.NodeLocal{
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

func NewClient(conf *config.GeoIP) (*Client, error) {
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
		return nil, fmt.Errorf("run geoipupdate: %w", err)
	}

	reader, err := geoip2.Open(conf.File)
	if err != nil {
		return nil, fmt.Errorf("open geolite2: %w", err)
	}

	return &Client{
		reader: reader,
	}, nil
}
