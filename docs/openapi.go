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
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/activity"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/rss3-network/protocol-go/schema/typex"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
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
	platforms := worker.PlatformStrings()

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
	schemas := make(map[string]map[string]interface{})

	schemas[tag.Transaction.String()] = generateTransactionMetadataActionSchemas()
	schemas[tag.Collectible.String()] = generateCollectibleMetadataActionSchemas()
	schemas[tag.Exchange.String()] = generateExchangeMetadataActionSchemas()
	schemas[tag.Social.String()] = generateSocialMetadataActionSchemas()
	schemas[tag.Metaverse.String()] = generateMetaverseMetadataActionSchemas()
	schemas[tag.RSS.String()] = generateRSSMetadataActionSchemas()

	var err error

	// Add the generated schemas to the components.schemas section of the OpenAPI document
	for tagx := range schemas {
		for metadatax, value := range schemas[tagx] {
			key := fmt.Sprintf("%s%s", tagx, metadatax)

			file, err = sjson.SetBytes(file, fmt.Sprintf("components.schemas.%s", key), value)
			if err != nil {
				return nil, fmt.Errorf("sjson set schema err: %w", err)
			}
		}
	}

	// Prepare anyOf array for components.schemas.Action.properties.metadata
	anyOfArray := make([]map[string]interface{}, 0)

	for key := range schemas {
		anyOfArray = append(anyOfArray, map[string]interface{}{
			"$ref": fmt.Sprintf("#/components/schemas/%s", key),
		})
	}

	// Add the anyOf array to components.schemas.Action.properties.metadata
	file, err = sjson.SetBytes(file, "components.schemas.Action.properties.metadata.anyOf", anyOfArray)
	if err != nil {
		return nil, fmt.Errorf("sjson set anyOf err: %w", err)
	}

	return file, nil
}

func generateTransactionMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.TransactionApproval.String(): generateMetadataObject(reflect.TypeOf(metadata.TransactionApproval{})),
		typex.TransactionBridge.String():   generateMetadataObject(reflect.TypeOf(metadata.TransactionBridge{})),
		typex.TransactionTransfer.String(): generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
		typex.TransactionBurn.String():     generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
		typex.TransactionMint.String():     generateMetadataObject(reflect.TypeOf(metadata.TransactionTransfer{})),
	}
}

func generateCollectibleMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.CollectibleApproval.String(): generateMetadataObject(reflect.TypeOf(metadata.CollectibleApproval{})),
		typex.CollectibleTrade.String():    generateMetadataObject(reflect.TypeOf(metadata.CollectibleTrade{})),
		typex.CollectibleTransfer.String(): generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
		typex.CollectibleBurn.String():     generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
		typex.CollectibleMint.String():     generateMetadataObject(reflect.TypeOf(metadata.CollectibleTransfer{})),
	}
}

func generateExchangeMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.ExchangeLiquidity.String(): generateMetadataObject(reflect.TypeOf(metadata.ExchangeLiquidity{})),
		typex.ExchangeStaking.String():   generateMetadataObject(reflect.TypeOf(metadata.ExchangeStaking{})),
		typex.ExchangeSwap.String():      generateMetadataObject(reflect.TypeOf(metadata.ExchangeSwap{})),
	}
}

func generateSocialMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.SocialPost.String():    generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialComment.String(): generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialRevise.String():  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialReward.String():  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialShare.String():   generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialDelete.String():  generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialMint.String():    generateMetadataObject(reflect.TypeOf(metadata.SocialPost{})),
		typex.SocialProfile.String(): generateMetadataObject(reflect.TypeOf(metadata.SocialProfile{})),
		typex.SocialProxy.String():   generateMetadataObject(reflect.TypeOf(metadata.SocialProxy{})),
	}
}

func generateMetaverseMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.MetaverseBurn.String():     generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseMint.String():     generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseTransfer.String(): generateMetadataObject(reflect.TypeOf(metadata.MetaverseTransfer{})),
		typex.MetaverseTrade.String():    generateMetadataObject(reflect.TypeOf(metadata.MetaverseTrade{})),
	}
}

func generateRSSMetadataActionSchemas() map[string]interface{} {
	return map[string]interface{}{
		typex.RSSFeed.String(): generateMetadataObject(reflect.TypeOf(metadata.RSS{})),
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

				if elemType == reflect.TypeOf(decimal.Decimal{}) || elemType == reflect.TypeOf(big.Int{}) {
					properties[fieldName] = map[string]interface{}{
						"type": "string",
					}

					continue
				}

				if elemType.Kind() == reflect.Struct {
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

		if elemType == reflect.TypeOf(decimal.Decimal{}) || elemType == reflect.TypeOf(big.Int{}) {
			return "string"
		}

		return transformOpenAPIType(elemType)
	case reflect.Struct:
		if t == reflect.TypeOf(decimal.Decimal{}) || t == reflect.TypeOf(big.Int{}) {
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
