package filter

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-node/schema/filter"
)

func PlatformTransferWorker(platform filter.Platform) (string, error) {
	switch platform {
	case filter.PlatformRSS3:
		return filter.RSS3.String(), nil
	case filter.PlatformMirror:
		return filter.Mirror.String(), nil
	case filter.PlatformFarcaster:
		return filter.Farcaster.String(), nil
	case filter.PlatformParagraph:
		return filter.Paragraph.String(), nil
	case filter.PlatformOpenSea:
		return filter.OpenSea.String(), nil
	case filter.PlatformUniswap:
		return filter.Uniswap.String(), nil
	case filter.PlatformOptimism:
		return filter.Optimism.String(), nil
	case filter.PlatformAavegotchi:
		return filter.Aavegotchi.String(), nil
	case filter.PlatformLens:
		return filter.Lens.String(), nil
	default:
		return "", fmt.Errorf("%s belongs unknown worker", platform.String())
	}
}

func TagTransferWorker(tag filter.Tag) ([]string, error) {
	switch tag {
	case filter.TagTransaction:
		return []string{filter.Optimism.String()}, nil
	case filter.TagCollectible:
		return []string{filter.OpenSea.String()}, nil
	case filter.TagExchange:
		return []string{filter.RSS3.String(), filter.Uniswap.String()}, nil
	case filter.TagSocial:
		return []string{filter.Farcaster.String(), filter.Mirror.String(), filter.Lens.String(), filter.Paragraph.String()}, nil
	case filter.TagMetaverse:
		return []string{filter.Aavegotchi.String()}, nil
	default:
		return nil, fmt.Errorf("%s belongs unknown worker", tag.String())
	}
}
