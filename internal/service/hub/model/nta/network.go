package nta

type NetworkRequest struct {
	Network string `param:"network" validate:"required"`
}
