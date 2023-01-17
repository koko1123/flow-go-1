package rest

import (
	"github.com/koko1123/flow-go-1/access"
	"github.com/koko1123/flow-go-1/engine/access/rest/models"
	"github.com/koko1123/flow-go-1/engine/access/rest/request"
)

// GetNetworkParameters returns network-wide parameters of the blockchain
func GetNetworkParameters(r *request.Request, backend access.API, link models.LinkGenerator) (interface{}, error) {
	params := backend.GetNetworkParameters(r.Context())

	var response models.NetworkParameters
	response.Build(&params)
	return response, nil
}
