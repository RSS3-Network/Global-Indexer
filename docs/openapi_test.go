package docs

import (
	"reflect"
	"testing"

	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMetadataObject(t *testing.T) {
	t.Run("embedding of struct", func(t *testing.T) {
		t.Parallel()

		result := generateMetadataObject(reflect.TypeOf(metadata.MetaverseTrade{}))

		assert.Equal(t, "object", result["type"])
		assert.Contains(t, result["properties"], "action")
		assert.Contains(t, result["properties"], "cost")
		assert.Contains(t, result["properties"], "address")
		assert.NotContains(t, result["properties"], "Token")
	})

	t.Run("struct array", func(t *testing.T) {
		t.Parallel()

		result := generateMetadataObject(reflect.TypeOf(metadata.ExchangeLiquidity{}))

		assert.Equal(t, "object", result["type"])
		assert.Contains(t, result["properties"], "tokens")
	})
}
