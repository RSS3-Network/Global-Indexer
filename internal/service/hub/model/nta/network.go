package nta

type NetworkRequest struct {
	NodeRequest

	NetworkName string `param:"network_name" validate:"required"`
}

type WorkerRequest struct {
	NetworkRequest

	WorkerName string `param:"worker_name" validate:"required"`
}
