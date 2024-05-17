package nta

type NetworkRequest struct {
	NetworkName string `param:"network_name" validate:"required"`
}

type WorkerRequest struct {
	NetworkRequest

	WorkerName string `param:"worker_name" validate:"required"`
}
