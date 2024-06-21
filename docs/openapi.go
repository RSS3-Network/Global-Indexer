package docs

import (
	"embed"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/activity"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/rss3-network/protocol-go/schema/typex"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed openapi.json
var EmbedFS embed.FS

var FilePath = "./docs/openapi.json"

func Generate() ([]byte, error) {
	// Read the existing openapi.json file from EmbedFS
	file, err := EmbedFS.ReadFile("openapi.json")
	if err != nil {
		zap.L().Error("read file error", zap.Error(err))
		return nil, err
	}

	// Generate network, tag, platform, direction enum.
	if file, err = generateEnum(file); err != nil {
		return nil, err
	}

	// Generate metadata action schema.
	if file, err = generateMetadataAction(file); err != nil {
		return nil, err
	}

	// Write the updated file back to the file system
	_ = os.WriteFile(FilePath, file, 0600)

	return file, nil
}

func generateEnum(file []byte) ([]byte, error) {
	var err error

	// Generate error code values.
	file, err = sjson.SetBytes(file, "components.schemas.ResponseError.properties.error_code.enum", errorx.ErrorCodeStrings())
	if err != nil {
		return nil, fmt.Errorf("sjson set error code enum err: %w", err)
	}

	// Generate network values.
	networks := lo.Filter(network.NetworkStrings(), func(s string, _ int) bool {
		return !lo.Contains([]string{
			network.Unknown.String(),
			network.Bitcoin.String(),
			network.SatoshiVM.String(),
			network.RSS.String(),
		}, s)
	})

	sort.Strings(networks)

	file, err = sjson.SetBytes(file, "components.schemas.Network.enum", networks)
	if err != nil {
		return nil, fmt.Errorf("sjson set network enum err: %w", err)
	}

	// Generate tag values.
	tags := tag.TagStrings()

	sort.Strings(tags)

	file, err = sjson.SetBytes(file, "components.schemas.Tag.enum", tags)
	if err != nil {
		return nil, fmt.Errorf("sjson set tag enum err: %w", err)
	}

	// Generate platform values.
	platforms := decentralized.PlatformStrings()

	sort.Strings(platforms)

	file, err = sjson.SetBytes(file, "components.schemas.Platform.enum", platforms)
	if err != nil {
		return nil, fmt.Errorf("sjson set platform enum err: %w", err)
	}

	// Generate direction values.
	file, err = sjson.SetBytes(file, "components.schemas.Direction.enum", activity.DirectionStrings())
	if err != nil {
		return nil, fmt.Errorf("sjson set direction enum err: %w", err)
	}

	// Generate type values.
	types := make([]string, 0)

	for _, v := range tag.TagValues() {
		for _, t := range schema.GetTypesByTag(v) {
			types = append(types, t.Name())
		}
	}

	types = lo.Uniq(types)

	sort.Strings(types)

	file, err = sjson.SetBytes(file, "components.schemas.Type.enum", types)
	if err != nil {
		return nil, fmt.Errorf("sjson set type enum err: %w", err)
	}

	return file, nil
}

func generateMetadataAction(file []byte) ([]byte, error) {
	// Generate all MetadataActionSchemas
	schemas := make(map[tag.Tag]map[schema.Type]interface{})

	schemas[tag.Transaction] = generateTransactionMetadataActionSchemas()
	schemas[tag.Collectible] = generateCollectibleMetadataActionSchemas()
	schemas[tag.Exchange] = generateExchangeMetadataActionSchemas()
	schemas[tag.Social] = generateSocialMetadataActionSchemas()
	schemas[tag.Metaverse] = generateMetaverseMetadataActionSchemas()
	schemas[tag.RSS] = generateRSSMetadataActionSchemas()

	anyOfArray := make([]map[string]interface{}, 0)

	var err error

	// Add the generated schemas to the components.schemas section of the OpenAPI document
	for tagx := range schemas {
		for typex, value := range schemas[tagx] {
			key := fmt.Sprintf("%s%s", toTitleCase(tagx.String()), toTitleCase(typex.Name()))

			file, err = sjson.SetBytes(file, fmt.Sprintf("components.schemas.%s", key), value)
			if err != nil {
				return nil, fmt.Errorf("sjson set schema err: %w", err)
			}

			// Prepare anyOf array for components.schemas.Action.properties.metadata
			anyOfArray = append(anyOfArray, map[string]interface{}{
				"$ref": fmt.Sprintf("#/components/schemas/%s", key),
			})
		}
	}

	// Add the anyOf array to components.schemas.Action.properties.metadata
	file, err = sjson.SetBytes(file, "components.schemas.Action.properties.metadata.anyOf", anyOfArray)
	if err != nil {
		return nil, fmt.Errorf("sjson set anyOf err: %w", err)
	}

	return file, nil
}

// toTitleCase converts a string to title case
func toTitleCase(s string) string {
	caser := cases.Title(language.English, cases.NoLower)

	return caser.String(s)
}

func generateTransactionMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.TransactionApproval: generateMetadataObject(reflect.TypeOf(metadata.TransactionApproval{})),
		typex.TransactionBridge:   generateMetadataObject(reflect.TypeOf(metadata.TransactionBridge{})),
		typex.TransactionTransfer: generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
		typex.TransactionBurn:     generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
		typex.TransactionMint:     generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
	}
}

func generateCollectibleMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.CollectibleApproval: generateMetadataObject(reflect.TypeOf(metadata.CollectibleApproval{})),
		typex.CollectibleTrade:    generateMetadataObject(reflect.TypeOf(metadata.CollectibleTrade{})),
		typex.CollectibleTransfer: generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
		typex.CollectibleBurn:     generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
		typex.CollectibleMint:     generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
	}
}

func generateExchangeMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.ExchangeLiquidity: generateMetadataObject(reflect.TypeOf(metadata.ExchangeLiquidity{})),
		typex.ExchangeStaking:   generateMetadataObject(reflect.TypeOf(metadata.ExchangeStaking{})),
		typex.ExchangeSwap:      generateMetadataObject(reflect.TypeOf(metadata.ExchangeSwap{})),
	}
}

func generateSocialMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.SocialPost:    generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialComment: generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialRevise:  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialReward:  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialShare:   generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialDelete:  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialMint:    generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialProfile: generateMetadataObject(reflect.TypeOf(metadata.SocialProfile{})),
		typex.SocialProxy:   generateMetadataObject(reflect.TypeOf(metadata.SocialProxy{})),
	}
}

func generateMetaverseMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.MetaverseBurn:     generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseMint:     generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseTransfer: generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseTrade:    generateMetadataObject(reflect.TypeOf(metadata.MetaverseTrade{})),
	}
}

func generateRSSMetadataActionSchemas() map[schema.Type]interface{} {
	return map[schema.Type]interface{}{
		typex.RSSFeed: generateMetadataObject(reflect.TypeOf(metadata.RSS{})),
	}
}

func generateMetadataObject(t reflect.Type) map[string]interface{} {
	object := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	properties := object["properties"].(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Remove 'omitempty' from the JSON tag to get the actual field name
		fieldName := strings.Split(field.Tag.Get("json"), ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// Check if the field type has a corresponding Strings function
		if method, ok := hasEnumStringsFunction(field.Type); ok {
			fieldSchema := map[string]interface{}{
				"type": "string",
				"enum": getEnumStrings(method),
			}
			properties[fieldName] = fieldSchema
		} else {
			// Handle pointer types first
			if field.Type.Kind() == reflect.Ptr {
				elemType := field.Type.Elem()

				if elemType == reflect.TypeOf(big.Int{}) {
					properties[fieldName] = map[string]interface{}{
						"type": "string",
					}

					continue
				}

				if elemType.Kind() == reflect.Struct {
					if elemType == reflect.TypeOf(metadata.SocialPost{}) && fieldName == "target" {
						properties[fieldName] = map[string]interface{}{
							"$ref": "#/components/schemas/SocialPost",
						}

						continue
					}

					properties[fieldName] = generateMetadataObject(elemType)

					continue
				}
			}

			// Then handle struct types
			if field.Type.Kind() == reflect.Struct {
				fieldSchema := generateMetadataObject(field.Type)
				properties[fieldName] = fieldSchema
			} else {
				fieldSchema := map[string]interface{}{
					"type": transformOpenAPIType(field.Type),
				}
				properties[fieldName] = fieldSchema
			}
		}
	}

	return object
}

func transformOpenAPIType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Ptr:
		elemType := t.Elem()

		if elemType == reflect.TypeOf(big.Int{}) {
			return "string"
		}

		return transformOpenAPIType(elemType)
	case reflect.Struct:
		if t == reflect.TypeOf(big.Int{}) {
			return "string"
		}

		return "object"
	default:
		return "string"
	}
}

func hasEnumStringsFunction(t reflect.Type) (reflect.Value, bool) {
	if t.Name() == "" {
		return reflect.Value{}, false
	}

	funcName := t.Name() + "Strings"
	globalFuncs := []interface{}{
		metadata.TransactionApprovalActionStrings,
		metadata.TransactionBridgeActionStrings,
		metadata.ExchangeLiquidityActionStrings,
		metadata.ExchangeStakingActionStrings,
		metadata.SocialProfileActionStrings,
		metadata.SocialProfileActionStrings,
		metadata.MetaverseTradeActionStrings,
		metadata.StandardStrings,
		network.NetworkStrings,
	}

	for _, f := range globalFuncs {
		funcFullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.HasSuffix(funcFullName, funcName) {
			return reflect.ValueOf(f), true
		}
	}

	return reflect.Value{}, false
}

func getEnumStrings(method reflect.Value) []string {
	results := method.Call(nil)
	return results[0].Interface().([]string)
}
