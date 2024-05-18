package nta

type GetNetworkRequest struct {
	NetworkName string `param:"network_name" validate:"required"`
}

type GetWorkerRequest struct {
	NetworkName string `param:"network_name" validate:"required"`
	WorkerName  string `param:"worker_name" validate:"required"`
}
