package docs

import (
	"reflect"
	"testing"
)

type Metadata struct {
	Token
}

type Token struct {
	Address        *string `json:"address,omitempty"`
	Name           string  `json:"name,omitempty"`
	Symbol         string  `json:"symbol,omitempty"`
	URI            string  `json:"uri,omitempty"`
	Decimals       uint8   `json:"decimals,omitempty"`
	ParsedImageURL string  `json:"parsed_image_url,omitempty"`
}

func TestGenerateMetadataObject(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		result := generateMetadataObject(reflect.TypeOf(Metadata{}))

		t.Log(result)
	})
}
