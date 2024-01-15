package hub

const (
	PathGetNodes      = "/nodes"
	PathGetNode       = "/node/:id"
	PathNodesRegister = "/nodes/register"
	PathStaking       = "/staking"
	PathBridging      = "/bridging"

	// http://localhost/rss/bilibili/user/coin/5374948?limit=1
	PathGetRSSHub = "/rss/*"
	// http://localhost/decentralized/tx/0x0000000000000000000000005cfabab80f9b95b55de329eecbd827c95ea5c814
	PathGetDecentralizedTx = "/decentralized/tx/:id"
	// http://localhost/decentralized/0x21f73438F3B1b00374054725A1C0D606831ECdA9
	PathGetDecentralizedActivities = "/decentralized/:account"
)
