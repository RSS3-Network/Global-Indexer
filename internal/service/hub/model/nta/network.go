package nta

type NetworkRequest struct {
	NetworkName string `param:"network_name" validate:"required"`
}

type WorkerRequest struct {
	NetworkRequest

	WorkerName string `param:"worker_name" validate:"required"`
}

// NetworkParamsData contains the network parameters
type NetworkParamsData struct {
	NetworkAssets map[string]Asset `json:"network_assets"`
	WorkerAssets  map[string]Asset `json:"worker_assets"`
}

type Asset struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Platform string `json:"platform,omitempty"`
	IconURL  string `json:"icon_url"`
}
