package nta

type NetworkRequest struct {
	Network string `param:"network" validate:"required"`
}

type WorkerRequest struct {
	NetworkRequest

	Worker string `param:"worker" validate:"required"`
}
