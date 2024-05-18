package hub

import (
	"context"
	"fmt"
	"go/types"
	"net"
	"net/http"

	"github.com/a-h/rest"
	"github.com/a-h/rest/swaggerui"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/docs"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/constant"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
)

const Name = "hub"

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = "80"
)

type Server struct {
	httpServer *echo.Echo
	hub        *Hub
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) Run(_ context.Context) error {
	address := net.JoinHostPort(DefaultHost, DefaultPort)

	return s.httpServer.Start(address)
}

func NewServer(databaseClient database.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, ethereumMultiChainClient *ethereum.MultiChainClient, nameService *nameresolver.NameResolver, httpClient httputil.Client) (service.Server, error) {
	hub, err := NewHub(context.Background(), databaseClient, redisClient, ethereumMultiChainClient, geoLite2, nameService, httpClient)
	if err != nil {
		return nil, fmt.Errorf("new hub: %w", err)
	}

	instance := Server{
		httpServer: echo.New(),
		hub:        hub,
	}

	// init restAPI for generating swagger ui
	restAPI := rest.NewAPI(constant.BuildServiceName())
	restAPI.StripPkgPaths = []string{
		"github.com/rss3-network",
		"math",
		"github.com/shopspring/decimal",
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Validator = defaultValidator
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	{
		groupName := "/nodes"
		group := instance.httpServer.Group(groupName)

		// Node registration and heartbeat
		// nodes.POST("/register", instance.hub.nta.RegisterNode)
		SetupRoute[nta.RegisterNodeRequest, string](group, restAPI, docs.RouteOption{
			Method:     http.MethodPost,
			Path:       "/register",
			GroupNames: []string{groupName},
			Handler:    instance.hub.nta.RegisterNode,
		})
		// nodes.POST("/heartbeat", instance.hub.nta.NodeHeartbeat)
		SetupRoute[nta.NodeHeartbeatRequest, string](group, restAPI, docs.RouteOption{
			Method:     http.MethodPost,
			Path:       "/heartbeat",
			GroupNames: []string{groupName},
			Handler:    instance.hub.nta.NodeHeartbeat,
		})
	}

	// nta is short for Network Transparency API
	ntaGroupName := "/nta"
	ntaGroup := instance.httpServer.Group(ntaGroupName)
	{
		{
			groupName := "/bridge"
			group := ntaGroup.Group(groupName)
			// bridge.GET("/transactions", instance.hub.nta.GetBridgeTransactions)
			SetupRoute[nta.GetBridgeTransactionsRequest, nta.TypedResponse[nta.GetBridgeTransactionsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/transactions",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetBridgeTransactions,
			})
			// bridge.GET("/transactions/:transaction_hash", instance.hub.nta.GetBridgeTransaction)
			SetupRoute[nta.GetBridgeTransactionRequest, nta.TypedResponse[nta.BridgeTransaction]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/transactions/:transaction_hash",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetBridgeTransaction,
			})
		}

		{
			groupName := "/chips"
			group := ntaGroup.Group(groupName)
			// chips.GET("", instance.hub.nta.GetStakeChips)
			SetupRoute[nta.GetStakeChipsRequest, nta.TypedResponse[nta.GetStakeChipsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeChips,
			})

			// chips.GET("/:chip_id", instance.hub.nta.GetStakeChip)
			SetupRoute[nta.GetStakeChipRequest, nta.TypedResponse[nta.GetStakeChipResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:chip_id",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeChip,
			})

			// chips.GET("/:chip_id/image.svg", instance.hub.nta.GetStakeChipImage)
			// FIXME: only application/json media type is supported by rest.API
			SetupRoute[nta.GetStakeChipRequest, []byte](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:chip_id/image.svg",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeChipImage,
			})
		}

		{
			groupName := "/epochs"
			group := ntaGroup.Group(groupName)

			// epochs.GET("", instance.hub.nta.GetEpochs)
			SetupRoute[nta.GetEpochsRequest, nta.TypedResponse[nta.GetEpochResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetEpochs,
			})

			// epochs.GET("/:epoch_id", instance.hub.nta.GetEpoch)
			SetupRoute[nta.GetEpochRequest, nta.TypedResponse[nta.GetEpochsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:epoch_id",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetEpoch,
			})

			// epochs.GET("/:node_address/rewards", instance.hub.nta.GetEpochNodeRewards)
			SetupRoute[nta.GetEpochNodeRewardsRequest, nta.TypedResponse[nta.GetEpochsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address/rewards",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetEpochNodeRewards,
			})

			// epochs.GET("/distributions/:transaction_hash", instance.hub.nta.GetEpochDistribution)
			SetupRoute[nta.GetEpochDistributionRequest, nta.TypedResponse[nta.GetEpochDistributionResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/distributions/:transaction_hash",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetEpochDistribution,
			})
		}

		{
			groupName := "/networks"
			group := ntaGroup.Group(groupName)

			// FIXME: the response it untyped
			// networks.GET("", instance.hub.nta.GetNetworks)
			SetupRoute[types.Nil, string](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNetworks,
			})

			// FIXME: the response it untyped
			// group.GET("/:network_name/list_workers", instance.hub.nta.GetWorkersByNetwork)
			SetupRoute[nta.GetNetworkRequest, string](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:network_name/list_workers",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetWorkersByNetwork,
			})

			// FIXME: the response it untyped
			// group.GET("/:network_name/workers/:worker_name", instance.hub.nta.GetWorkerDetail)
			SetupRoute[nta.GetWorkerRequest, string](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:network_name/workers/:worker_name",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetWorkerDetail,
			})
		}

		{
			groupName := "/nodes"
			group := ntaGroup.Group(groupName)

			// group.GET("", instance.hub.nta.GetNodes)
			SetupRoute[nta.BatchNodeRequest, nta.TypedResponse[nta.GetNodesResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodes,
			})

			// group.GET("/:node_address", instance.hub.nta.GetNode)
			SetupRoute[nta.GetNodeRequest, nta.TypedResponse[nta.GetNodeResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNode,
			})

			// group.GET("/:node_address/avatar.svg", instance.hub.nta.GetNodeAvatar)
			// FIXME: only application/json media type is supported by rest.API
			SetupRoute[nta.GetNodeRequest, nta.TypedResponse[nta.GetNodeResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address/avatar.svg",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeAvatar,
			})

			// group.GET("/:node_address/challenge", instance.hub.nta.GetNodeChallenge)
			SetupRoute[nta.GetNodeChallengeRequest, nta.TypedResponse[nta.GetNodeChallengeResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address/challenge",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeChallenge,
			})

			// group.GET("/:node_address/events", instance.hub.nta.GetNodeEvents)
			SetupRoute[nta.GetNodeEventsRequest, nta.TypedResponse[nta.GetNodeEventsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address/events",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeEvents,
			})

			// group.POST("/:node_address/hide_tax_rate", instance.hub.nta.PostNodeHideTaxRate)
			// hide for now
			// SetupRoute[nta.PostNodeHideTaxRateRequest, types.Nil](group, restAPI, docs.RouteOption{
			//	Method:     http.MethodPost,
			//	Path:       "/:node_address/hide_tax_rate",
			//	GroupNames: []string{ntaGroupName, groupName},
			//	Handler:    instance.hub.nta.PostNodeHideTaxRate,
			// })

			// group.GET("/:node_address/operation/profit", instance.hub.nta.GetNodeOperationProfit)
			SetupRoute[nta.GetNodeOperationProfitRequest, nta.TypedResponse[nta.GetNodeOperationProfitResponse]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:node_address/operation/profit",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeOperationProfit,
			})
		}

		{
			groupName := "/snapshots"
			group := ntaGroup.Group(groupName)
			// group.GET("/nodes/count", instance.hub.nta.GetNodeCountSnapshots)
			SetupRoute[types.Nil, nta.TypedResponse[nta.GetNodeCountSnapshotsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/nodes/count",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeCountSnapshots,
			})

			// group.GET("/nodes/operation/profit", instance.hub.nta.GetNodeOperationProfitSnapshots)
			SetupRoute[nta.GetNodeOperationProfitSnapshotsRequest, nta.TypedResponse[[]schema.OperatorProfitSnapshot]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/nodes/operation/profit",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetNodeOperationProfitSnapshots,
			})

			// group.GET("/stakers/count", instance.hub.nta.GetStakerCountSnapshots)
			SetupRoute[types.Nil, nta.TypedResponse[nta.GetStakerProfitSnapshotsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/stakers/count",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakerCountSnapshots,
			})

			// group.GET("/stakers/profit", instance.hub.nta.GetStakerProfitSnapshots)
			SetupRoute[nta.GetStakerProfitSnapshotsRequest, nta.TypedResponse[[]schema.StakerProfitSnapshot]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/stakers/profit",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakerProfitSnapshots,
			})

			// group.POST("/nodes/min_tokens_to_stake", instance.hub.nta.BatchGetNodeMinTokensToStakeSnapshots)
			SetupRoute[nta.BatchNodeMinTokensToStakeRequest, nta.TypedResponse[nta.BatchGetNodeMinTokensToStakeSnapshotsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodPost,
				Path:       "/nodes/min_tokens_to_stake",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.BatchGetNodeMinTokensToStakeSnapshots,
			})
		}

		{
			groupName := "/stake"
			group := ntaGroup.Group(groupName)
			// group.GET("/:staker_address/profit", instance.hub.nta.GetStakeOwnerProfit)
			SetupRoute[nta.GetStakeOwnerProfitRequest, nta.TypedResponse[nta.GetStakeOwnerProfitResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/:staker_address/profit",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeOwnerProfit,
			})

			// group.GET("/stakings", instance.hub.nta.GetStakeStakings)
			SetupRoute[nta.GetStakeStakingsRequest, nta.TypedResponse[nta.GetStakeStakingsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/stakings",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeStakings,
			})

			// group.GET("/transactions", instance.hub.nta.GetStakeTransactions)
			SetupRoute[nta.GetStakeTransactionsRequest, nta.TypedResponse[nta.GetStakeTransactionsResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/transactions",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeTransactions,
			})

			// group.GET("/transactions/:transaction_hash", instance.hub.nta.GetStakeTransaction)
			SetupRoute[nta.GetStakeTransactionRequest, nta.TypedResponse[nta.GetStakeTransactionResponseData]](group, restAPI, docs.RouteOption{
				Method:     http.MethodGet,
				Path:       "/transactions/:transaction_hash",
				GroupNames: []string{ntaGroupName, groupName},
				Handler:    instance.hub.nta.GetStakeTransaction,
			})
		}
	}

	dslGroupName := "/dsl"
	dslGroup := instance.httpServer.Group(dslGroupName)
	{
		groupName := "/rss"
		group := dslGroup.Group(groupName)

		// group.GET("/*", instance.hub.dslGroup.GetRSSHub)
		// FIXME: the response it untyped
		SetupRoute[types.Nil, []byte](group, restAPI, docs.RouteOption{
			Method:     http.MethodGet,
			Path:       "",
			GroupNames: []string{dslGroupName, groupName},
			Handler:    instance.hub.dsl.GetRSSHub,
		})
	}

	{
		groupName := "/decentralized"
		group := dslGroup.Group(groupName)
		// group.GET("/tx/:id", instance.hub.dslGroup.GetActivity)
		// FIXME: the response it untyped
		SetupRoute[dsl.ActivityRequest, []byte](group, restAPI, docs.RouteOption{
			Method:     http.MethodGet,
			Path:       "/tx/:id",
			GroupNames: []string{dslGroupName, groupName},
			Handler:    instance.hub.dsl.GetActivity,
		})
		// group.GET("/:account", instance.hub.dsl.GetAccountActivities)
		SetupRoute[dsl.ActivitiesRequest, []byte](group, restAPI, docs.RouteOption{
			Method:     http.MethodGet,
			Path:       "/:account",
			GroupNames: []string{dslGroupName, groupName},
			Handler:    instance.hub.dsl.GetAccountActivities,
		})
	}

	{
		// create OpenAPI spec
		spec, err := restAPI.Spec()
		if err != nil {
			return nil, fmt.Errorf("failed to create spec: %w", err)
		}

		spec.Info.Version = constant.BuildServiceName()
		spec.Info.Description = constant.ServiceName

		// attach the UI Handler.
		ui, err := swaggerui.New(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to create swagger UI Handler: %w", err)
		}

		instance.httpServer.GET("/swagger-ui*", echo.WrapHandler(ui))
	}

	return &instance, nil
}

// SetupRoute is a helper function to setup a route in both echo router and rest.API
func SetupRoute[Request, Response any](router *echo.Group, api *rest.API, routeOption docs.RouteOption) {
	route := &rest.Route{}

	// we only need GET and POST for now
	switch routeOption.Method {
	case http.MethodGet:
		router.GET(routeOption.Path, routeOption.Handler)
		route = api.Get(routeOption.GetCompletePath())

	case http.MethodPost:
		router.POST(routeOption.Path, routeOption.Handler)
		route = api.Post(routeOption.GetCompletePath())
		route.HasRequestModel(rest.ModelOf[Request]())
	}

	params := docs.ConvertStructToParams[Request]()

	// sets the path parameters
	for name, param := range params.Path {
		route.HasPathParameter(name, param)
	}

	// sets the query parameters
	for name, param := range params.Query {
		route.HasQueryParameter(name, param)
	}

	// set the tags
	route.Tags = routeOption.GetTags()

	// set the response model
	route.
		HasResponseModel(http.StatusOK, rest.ModelOf[Response]())

	// set the error response models
	route.
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[errorx.ErrorResponse]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[errorx.ErrorResponse]())
}
