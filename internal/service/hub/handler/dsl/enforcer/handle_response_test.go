package enforcer

import (
	"errors"
	"testing"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/stretchr/testify/assert"
)

var (
	nullData                  = `{"data":null}`
	errResponseData           = `{"error":"it is error","error_code":"internal_error"}`
	activityResponseData0     = `{"data":{"id":"0xf6a9fa15432b27ac86f33a678635cbc97244cbe56457eb5dc0c946aec715c639","owner":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","network":"ethereum","index":0,"from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","tag":"transaction","type":"approval","fee":{"amount":"8126461890087948","decimal":18},"total_actions":2,"actions":[{"tag":"transaction","type":"transfer","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"23900000000000000000000","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"approval","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"action":"revoke","address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"0","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}}],"direction":"in","success":true,"timestamp":1710060911},"meta":{"totalPages":1}}`
	activityResponseData1     = `{"data":{"id":"0xf6a9fa15432b27ac86f33a678635cbc97244cbe56457eb5dc0c946aec715c63","owner":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","network":"ethereum","index":0,"from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","tag":"transaction","type":"approval","fee":{"amount":"8126461890087948","decimal":18},"total_actions":2,"actions":[{"tag":"transaction","type":"transfer","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"23900000000000000000000","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"approval","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"action":"revoke","address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"0","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}}],"direction":"in","success":true,"timestamp":1710060911},"meta":{"totalPages":1}}`
	activityResponseData2     = `{"data":{"id":"0xf6a9fa15432b27ac86f33a678635cbc97244cbe56457eb5dc0c946aec715c639","owner":"0x4cbab69108Aa72151EDa5A3c164eA86845f1843","network":"ethereum","index":0,"from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","tag":"transaction","type":"approval","fee":{"amount":"8126461890087948","decimal":18},"total_actions":2,"actions":[{"tag":"transaction","type":"transfer","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"23900000000000000000000","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"approval","from":"0x827431510a5D249cE4fdB7F00C83a3353F471848","to":"0x4cbab69108Aa72151EDa5A3c164eA86845f18438","metadata":{"action":"revoke","address":"0xc98D64DA73a6616c42117b582e832812e7B8D57F","value":"0","name":"RSS3","symbol":"RSS3","decimals":18,"standard":"ERC-20"}}],"direction":"in","success":true,"timestamp":1710060911},"meta":{"totalPages":1}}`
	activitiesResponseData    = `{"data":[{"id":"0x00003ef8175e535375ab4b2f8319488eb464496c140eb3d83b5e6a6509d06823","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xbEb5Fc579115071764c7423A4f12eDde41f106Ed","tag":"transaction","type":"bridge","platform":"Optimism","fee":{"amount":"4561943091807336","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"bridge","platform":"Optimism","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"action":"withdraw","source_network":"optimism","target_network":"ethereum","token":{"address":"0xDe30da39c46104798bB5aA3fe8B9e0e1F348163F","value":"3837000000000000000000","name":"Gitcoin","symbol":"GTC","decimals":18,"standard":"ERC-20"}}}],"direction":"out","success":true,"timestamp":1712265479},{"id":"0x0ffafa50658d97bddb5f203b82178665340426a6ae174c64eef683b8ca9c2409","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"390221892177518088","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711966335},{"id":"0x7c5f5d5478e651500490cf26adbba5ea1151ae62e68443cff30d1bc0491cf11d","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"46134958229847664","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711302631},{"id":"0x20e902cf88d9763e79da21020c4bff04e18438732b797d51fe3e74df61a435b4","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"97043494515138844","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711047924},{"id":"0x8660a0e7cbaafc9485dcd1325e9ec8c4e4fd4a0a5386dab69d26f6a94fc34bbf","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"524209298721134181","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1710839144},{"id":"0x286912082702d86f0e86e0a46e026c0ef7579d4b6e9ae60d516989184a00aa8f","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"37515642107389972","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1710705805},{"id":"0xbed1fcddc6405766a89160fdb631918f8b8e706fb6080dc220ad16e1bc67da54","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xdAC17F958D2ee523a2206206994597C13D831ec7","tag":"transaction","type":"transfer","fee":{"amount":"2511801986208783","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"transfer","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xFf3A3f9d7906f5635a3Ca0C830Cb4C6022De8aA0","metadata":{"address":"0xdAC17F958D2ee523a2206206994597C13D831ec7","value":"9983000000","name":"Tether USD","symbol":"USDT","decimals":6,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710538295},{"id":"0xa7e8da0fd2c82878e0334999836f436b97627bce8cbfc629dca4fd339e2e4518","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","tag":"transaction","type":"transfer","fee":{"amount":"7878450439415343","decimal":18},"total_actions":4,"actions":[{"tag":"transaction","type":"transfer","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xA86ACA6D7C393c06DcDC30473ea3D1b05c358DFF","metadata":{"address":"0x03AA6298F1370642642415EDC0db8b957783e8D6","value":"1226000000000000000000","name":"NetMind Token","symbol":"NMT","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0xA86ACA6D7C393c06DcDC30473ea3D1b05c358DFF","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","metadata":{"address":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48","value":"9984588553","name":"USD Coin","symbol":"USDC","decimals":6,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x3416cF6C708Da44DB2624D63ea0AAef7113527C6","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0xdAC17F958D2ee523a2206206994597C13D831ec7","value":"9983823155","name":"Tether USD","symbol":"USDT","decimals":6,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","to":"0x3416cF6C708Da44DB2624D63ea0AAef7113527C6","metadata":{"address":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48","value":"9984588553","name":"USD Coin","symbol":"USDC","decimals":6,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710538031},{"id":"0x321b4c528736dce4849b4fa541201cb798670ac2941e1a17e101dab85f9e8d09","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x2d8e09b546d0067acDB415329f0cB2204B198aA9","to":"0x0BbC02Ef7ce79A820B7EDda8d5D409ae7615e636","tag":"transaction","type":"burn","fee":{"amount":"7709876852394610","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"burn","from":"0x0000000000000000000000000000000000000000","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0x03AA6298F1370642642415EDC0db8b957783e8D6","value":"1226000000000000000000","name":"NetMind Token","symbol":"NMT","decimals":18,"standard":"ERC-20"}}],"direction":"in","success":true,"timestamp":1710538007},{"id":"0x10bea49a3e07ac1fb193190d48cdba8d4c97b9760ebc59704bad79ff77c05161","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","tag":"transaction","type":"transfer","fee":{"amount":"5128797405055820","decimal":18},"total_actions":2,"actions":[{"tag":"transaction","type":"transfer","from":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","to":"0x56dc588EF98701886855EcF823A85DD6FdCCC295","metadata":{"address":"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2","value":"80000000000000000","name":"Wrapped Ether","symbol":"WETH","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x56dc588EF98701886855EcF823A85DD6FdCCC295","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0x438e48ed4ce6beECF503D43b9dbD3C30d516e7FD","value":"109280518892380658743","name":"UWON","symbol":"UWON","decimals":18,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710533015}],"meta":{"cursor":"0x10bea49a3e07ac1fb193190d48cdba8d4c97b9760ebc59704bad79ff77c05161:ethereum"}}`
	mutableActivityResponse   = `{"data":{"id":"0x00000000000000000000000000000aadbebe59e6a53f0d552388f00f3cf8e5b1","owner":"0x0000000000000000000000000000000000000091","network":"farcaster","index":0,"from":"0x15426Ef0B2A15B3F7dcA513BD70F88A9481a0320","to":"0x4aF0919907ccdBFF6C80463A47B4Db24599CC0d5","tag":"social","type":"share","platform":"Farcaster","total_actions":2,"actions":[{"tag":"social","type":"share","platform":"Farcaster","from":"0xf0a68CD0e9AC293Ae4E9A5730852C9Cd0eaa56a0","to":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","metadata":{"handle":"tricksloaded","profile_id":"277353","publication_id":"0x00000aAdBebe59E6a53F0d552388f00f3cF8E5b1","target":{"handle":"terrytat","body":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","profile_id":"277460","publication_id":"0x6C14EAe3Ebda408B7d20968fB992b0DF7E8ec2cd"}}},{"tag":"social","type":"share","platform":"Farcaster","from":"0x0000000000000000000000000000000000000091","to":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","metadata":{"handle":"tricksloaded","profile_id":"277353","publication_id":"0x00000aAdBebe59E6a53F0d552388f00f3cF8E5b1","target":{"handle":"terrytat","body":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","profile_id":"277460","publication_id":"0x6C14EAe3Ebda408B7d20968fB992b0DF7E8ec2cd"}}}],"direction":"out","success":true,"timestamp":1711059936},"meta":{"totalPages":1}}`
	mutableActivitiesResponse = `{"data":[{"id":"0x00003ef8175e535375ab4b2f8319488eb464496c140eb3d83b5e6a6509d06823","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xbEb5Fc579115071764c7423A4f12eDde41f106Ed","tag":"transaction","type":"bridge","platform":"Optimism","fee":{"amount":"4561943091807336","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"bridge","platform":"Optimism","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"action":"withdraw","source_network":"optimism","target_network":"ethereum","token":{"address":"0xDe30da39c46104798bB5aA3fe8B9e0e1F348163F","value":"3837000000000000000000","name":"Gitcoin","symbol":"GTC","decimals":18,"standard":"ERC-20"}}}],"direction":"out","success":true,"timestamp":1712265479},{"id":"0x000000000000000000000000113f4b4c3765e5f05fd197c5c35b8a8a9b34245b","owner":"0xe5d6216F0085a7F6B9b692e06cf5856e6fA41B55","network":"farcaster","index":0,"from":"0xe5d6216F0085a7F6B9b692e06cf5856e6fA41B55","to":"0xe5d6216F0085a7F6B9b692e06cf5856e6fA41B55","tag":"social","type":"post","platform":"Farcaster","total_actions":1,"actions":[{"tag":"social","type":"post","platform":"Farcaster","from":"0x8888888198FbdC8c017870cC5d3c96D0cf15C4F0","to":"0x8888888198FbdC8c017870cC5d3c96D0cf15C4F0","metadata":{"handle":"brucexc.eth","body":"https://explorer.rss3.io/👍","media":[{"address":"https://explorer.rss3.io/","mime_type":"text/html; charset=utf-8"}],"profile_id":"14142","publication_id":"0x113f4B4C3765E5f05FD197C5c35b8a8a9b34245B"}}],"direction":"self","success":true,"timestamp":1710822312},{"id":"0x0ffafa50658d97bddb5f203b82178665340426a6ae174c64eef683b8ca9c2409","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"390221892177518088","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711966335},{"id":"0x7c5f5d5478e651500490cf26adbba5ea1151ae62e68443cff30d1bc0491cf11d","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"46134958229847664","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711302631},{"id":"0x20e902cf88d9763e79da21020c4bff04e18438732b797d51fe3e74df61a435b4","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"97043494515138844","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1711047924},{"id":"0x8660a0e7cbaafc9485dcd1325e9ec8c4e4fd4a0a5386dab69d26f6a94fc34bbf","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"524209298721134181","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1710839144},{"id":"0x286912082702d86f0e86e0a46e026c0ef7579d4b6e9ae60d516989184a00aa8f","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"polygon","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x111111125421cA6dc452d289314280a0f8842A65","tag":"unknown","type":"unknown","platform":"AAVE","fee":{"amount":"37515642107389972","decimal":18},"total_actions":0,"actions":[],"direction":"out","success":true,"timestamp":1710705805},{"id":"0xbed1fcddc6405766a89160fdb631918f8b8e706fb6080dc220ad16e1bc67da54","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xdAC17F958D2ee523a2206206994597C13D831ec7","tag":"transaction","type":"transfer","fee":{"amount":"2511801986208783","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"transfer","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xFf3A3f9d7906f5635a3Ca0C830Cb4C6022De8aA0","metadata":{"address":"0xdAC17F958D2ee523a2206206994597C13D831ec7","value":"9983000000","name":"Tether USD","symbol":"USDT","decimals":6,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710538295},{"id":"0xa7e8da0fd2c82878e0334999836f436b97627bce8cbfc629dca4fd339e2e4518","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","tag":"transaction","type":"transfer","fee":{"amount":"7878450439415343","decimal":18},"total_actions":4,"actions":[{"tag":"transaction","type":"transfer","from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0xA86ACA6D7C393c06DcDC30473ea3D1b05c358DFF","metadata":{"address":"0x03AA6298F1370642642415EDC0db8b957783e8D6","value":"1226000000000000000000","name":"NetMind Token","symbol":"NMT","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0xA86ACA6D7C393c06DcDC30473ea3D1b05c358DFF","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","metadata":{"address":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48","value":"9984588553","name":"USD Coin","symbol":"USDC","decimals":6,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x3416cF6C708Da44DB2624D63ea0AAef7113527C6","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0xdAC17F958D2ee523a2206206994597C13D831ec7","value":"9983823155","name":"Tether USD","symbol":"USDT","decimals":6,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","to":"0x3416cF6C708Da44DB2624D63ea0AAef7113527C6","metadata":{"address":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48","value":"9984588553","name":"USD Coin","symbol":"USDC","decimals":6,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710538031},{"id":"0x321b4c528736dce4849b4fa541201cb798670ac2941e1a17e101dab85f9e8d09","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x2d8e09b546d0067acDB415329f0cB2204B198aA9","to":"0x0BbC02Ef7ce79A820B7EDda8d5D409ae7615e636","tag":"transaction","type":"burn","fee":{"amount":"7709876852394610","decimal":18},"total_actions":1,"actions":[{"tag":"transaction","type":"burn","from":"0x0000000000000000000000000000000000000000","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0x03AA6298F1370642642415EDC0db8b957783e8D6","value":"1226000000000000000000","name":"NetMind Token","symbol":"NMT","decimals":18,"standard":"ERC-20"}}],"direction":"in","success":true,"timestamp":1710538007},{"id":"0x10bea49a3e07ac1fb193190d48cdba8d4c97b9760ebc59704bad79ff77c05161","owner":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","network":"ethereum","index":0,"from":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","to":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","tag":"transaction","type":"transfer","fee":{"amount":"5128797405055820","decimal":18},"total_actions":2,"actions":[{"tag":"transaction","type":"transfer","from":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","to":"0x56dc588EF98701886855EcF823A85DD6FdCCC295","metadata":{"address":"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2","value":"80000000000000000","name":"Wrapped Ether","symbol":"WETH","decimals":18,"standard":"ERC-20"}},{"tag":"transaction","type":"transfer","from":"0x56dc588EF98701886855EcF823A85DD6FdCCC295","to":"0x40B8Cef0Fa62aEF0050bF7D8bb62Cf065583B648","metadata":{"address":"0x438e48ed4ce6beECF503D43b9dbD3C30d516e7FD","value":"109280518892380658743","name":"UWON","symbol":"UWON","decimals":18,"standard":"ERC-20"}}],"direction":"out","success":true,"timestamp":1710533015}],"meta":{"cursor":"0x10bea49a3e07ac1fb193190d48cdba8d4c97b9760ebc59704bad79ff77c05161:ethereum"}}`
	activitiesResponseData0   = `{"data":[{"actions":[{"from":"0xAff3E9613A1fA8a0585e2d6bc47CE2eA3eDa221B","metadata":{"address":"0xe72575e0BF93d171D082B69EF11dE1F537086B15","id":"0","name":"SHY","standard":"ERC-1155","symbol":"Reward at usdtbonus.org","value":"1"},"related_urls":["https://blockscout.com/xdai/mainnet/tx/0x27edca80b812753c7ab963d6c81e03a7c0fa6cca2cc013c0a79e9ed0d38f089e"],"tag":"collectible","to":"0xbaBc85e78872F47532A57F13BF3576c5369f59Ff","type":"transfer"}],"calldata":{"function_hash":"0xc204642c"},"direction":"in","fee":{"amount":"1906159020509772","decimal":18},"from":"0xccaDd61056f17179846E4e5F6dF06619741EB069","id":"0x27edca80b812753c7ab963d6c81e03a7c0fa6cca2cc013c0a79e9ed0d38f089e","index":0,"network":"gnosis","owner":"0x23c46e912b34C09c4bCC97F4eD7cDd762cee408A","success":true,"tag":"collectible","timestamp":1714424720,"to":"0xe72575e0BF93d171D082B69EF11dE1F537086B15","total_actions":206,"type":"transfer"},{"actions":[{"from":"0xFC2d970A4e1E9464a1c065275F00347ed511e1d6","metadata":{"address":"0xe9BDa4e137704462D2d3495981625bF0925971F5","id":"0","name":"SHY","standard":"ERC-1155","symbol":"Reward at shibabonus.ac","value":"1"},"related_urls":["https://blockscout.com/xdai/mainnet/tx/0xcfc5365ef9d93bd3a9944198a037c98d8a5c49bb169885f5890188252a3cea5e"],"tag":"collectible","to":"0xbaBc85e78872F47532A57F13BF3576c5369f59Ff","type":"transfer"}],"calldata":{"function_hash":"0xc204642c"},"direction":"in","fee":{"amount":"1906159020509772","decimal":18},"from":"0x89AC40052B680A2bbDABf6b700b66cee555C3aB9","id":"0xcfc5365ef9d93bd3a9944198a037c98d8a5c49bb169885f5890188252a3cea5e","index":0,"network":"gnosis","owner":"0x23c46e912b34C09c4bCC97F4eD7cDd762cee408A","success":true,"tag":"collectible","timestamp":1714424720,"to":"0xe9BDa4e137704462D2d3495981625bF0925971F5","total_actions":206,"type":"transfer"}]}`
	activitiesResponseData1   = `{"data":[{"actions":[{"from":"0xFC2d970A4e1E9464a1c065275F00347ed511e1d6","metadata":{"address":"0xe9BDa4e137704462D2d3495981625bF0925971F5","id":"0","name":"SHY","standard":"ERC-1155","symbol":"Reward at shibabonus.ac","value":"1"},"related_urls":["https://blockscout.com/xdai/mainnet/tx/0xcfc5365ef9d93bd3a9944198a037c98d8a5c49bb169885f5890188252a3cea5e"],"tag":"collectible","to":"0xbaBc85e78872F47532A57F13BF3576c5369f59Ff","type":"transfer"}],"calldata":{"function_hash":"0xc204642c"},"direction":"in","fee":{"amount":"1906159020509772","decimal":18},"from":"0x89AC40052B680A2bbDABf6b700b66cee555C3aB9","id":"0xcfc5365ef9d93bd3a9944198a037c98d8a5c49bb169885f5890188252a3cea5e","index":0,"network":"gnosis","owner":"0x23c46e912b34C09c4bCC97F4eD7cDd762cee408A","success":true,"tag":"collectible","timestamp":1714424720,"to":"0xe9BDa4e137704462D2d3495981625bF0925971F5","total_actions":206,"type":"transfer"},{"actions":[{"from":"0xAff3E9613A1fA8a0585e2d6bc47CE2eA3eDa221B","metadata":{"address":"0xe72575e0BF93d171D082B69EF11dE1F537086B15","id":"0","name":"SHY","standard":"ERC-1155","symbol":"Reward at usdtbonus.org","value":"1"},"related_urls":["https://blockscout.com/xdai/mainnet/tx/0x27edca80b812753c7ab963d6c81e03a7c0fa6cca2cc013c0a79e9ed0d38f089e"],"tag":"collectible","to":"0xbaBc85e78872F47532A57F13BF3576c5369f59Ff","type":"transfer"}],"calldata":{"function_hash":"0xc204642c"},"direction":"in","fee":{"amount":"1906159020509772","decimal":18},"from":"0xccaDd61056f17179846E4e5F6dF06619741EB069","id":"0x27edca80b812753c7ab963d6c81e03a7c0fa6cca2cc013c0a79e9ed0d38f089e","index":0,"network":"gnosis","owner":"0x23c46e912b34C09c4bCC97F4eD7cDd762cee408A","success":true,"tag":"collectible","timestamp":1714424720,"to":"0xe72575e0BF93d171D082B69EF11dE1F537086B15","total_actions":206,"type":"transfer"}]}`
)

func TestCompareData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		src      []byte
		des      []byte
		expected bool
	}{
		{
			name:     "IdenticalActivityResponses",
			src:      []byte(activityResponseData0),
			des:      []byte(activityResponseData0),
			expected: true,
		},
		{
			name:     "IdenticalActivityResponsesWithMutableActivity",
			src:      []byte(mutableActivityResponse),
			des:      []byte(mutableActivityResponse),
			expected: true,
		},
		{
			name:     "DifferentActivityResponses",
			src:      []byte(activityResponseData0),
			des:      []byte(activityResponseData1),
			expected: false,
		},
		{
			name:     "IdenticalActivitiesResponses",
			src:      []byte(activitiesResponseData),
			des:      []byte(activitiesResponseData),
			expected: true,
		},
		{
			name:     "IdenticalActivitiesResponsesWithMutableActivities",
			src:      []byte(activitiesResponseData),
			des:      []byte(mutableActivitiesResponse),
			expected: true,
		},
		{
			name:     "DifferentActivitiesResponses",
			src:      []byte(activitiesResponseData),
			des:      []byte(nullData),
			expected: false,
		},
		{
			name:     "InvalidActivityResponse",
			src:      []byte(errResponseData),
			des:      []byte(activityResponseData0),
			expected: false,
		},
		{
			name:     "InvalidActivitiesResponse",
			src:      []byte(errResponseData),
			des:      []byte(activitiesResponseData),
			expected: false,
		},
		{
			name:     "IdenticalWithDifferentOrder",
			src:      []byte(activitiesResponseData0),
			des:      []byte(activitiesResponseData1),
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := isResponseIdentical(tc.src, tc.des)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUpdateRequestsBasedOnDataCompare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		responses       []*model.DataResponse
		requests        []int
		invalidRequests []int
	}{
		{
			name: "one_error_response",
			responses: []*model.DataResponse{
				{Err: errors.New(errResponseData)},
			},
			requests:        []int{0},
			invalidRequests: []int{1},
		},
		{
			name: "one_valid_response",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0), Valid: true},
			},
			requests:        []int{1},
			invalidRequests: []int{0},
		},
		{
			name: "two_error_responses",
			responses: []*model.DataResponse{
				{Err: errors.New("error1")},
				{Err: errors.New("error2")},
			},
			requests:        []int{0, 0},
			invalidRequests: []int{1, 1},
		},
		{
			name: "one_error_with_two_responses",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0), Valid: true},
				{Err: errors.New("error")},
			},
			requests:        []int{1, 0},
			invalidRequests: []int{0, 1},
		},
		{
			name: "two_responses_with_different_data",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0), Valid: true},
				{Data: []byte(activityResponseData1), Valid: true},
			},
			requests:        []int{1, 0},
			invalidRequests: []int{0, 0},
		},
		{
			name: "two_responses_with_same_data",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData1), Valid: true},
				{Data: []byte(activityResponseData1), Valid: true},
			},
			requests:        []int{2, 1},
			invalidRequests: []int{0, 0},
		},
		{
			name: "three_errors",
			responses: []*model.DataResponse{
				{Err: errors.New("error1")},
				{Err: errors.New("error2")},
				{Err: errors.New("error3")},
			},
			requests:        []int{0, 0, 0},
			invalidRequests: []int{1, 1, 1},
		},
		{
			name: "two_errors",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0), Valid: true},
				{Err: errors.New("error2")},
				{Err: errors.New("error3")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 1, 1},
		},
		{
			name: "one_error_with_same_data",
			responses: []*model.DataResponse{
				{Data: []byte(activitiesResponseData)},
				{Data: []byte(activitiesResponseData)},
				{Err: errors.New("error3")},
			},
			requests:        []int{2, 1, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "one_error_with_different_data",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData2)},
				{Err: errors.New("error3")},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "three_same_data",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData0)},
			},
			requests:        []int{2, 1, 1},
			invalidRequests: []int{0, 0, 0},
		},
		{
			name: "three_different_data",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData1)},
				{Data: []byte(activityResponseData2)},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 0},
		},
		{
			name: "two_different_data_01",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData2)},
			},
			requests:        []int{2, 1, 0},
			invalidRequests: []int{0, 0, 1},
		},
		{
			name: "two_different_data_02",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData1)},
				{Data: []byte(activityResponseData0)},
			},
			requests:        []int{2, 0, 1},
			invalidRequests: []int{0, 1, 0},
		},
		{
			name: "two_different_data_12_with_valid",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(activityResponseData1), Valid: true},
				{Data: []byte(activityResponseData1)},
			},
			requests:        []int{0, 1, 1},
			invalidRequests: []int{1, 0, 0},
		},
		{
			name: "two_different_data_12_with_invalid",
			responses: []*model.DataResponse{
				{Data: []byte(activityResponseData0)},
				{Data: []byte(nullData)},
				{Data: []byte(nullData)},
			},
			requests:        []int{1, 0, 0},
			invalidRequests: []int{0, 0, 0},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			updatePointsBasedOnIdentity(tc.responses)

			for i, result := range tc.responses {
				assert.Equal(t, tc.requests[i], result.ValidPoint)
				assert.Equal(t, tc.invalidRequests[i], result.InvalidPoint)
			}
		})
	}
}
