package id

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
)

// IdentityFilterIdentifierProvider implements an IdentifierProvider which provides the identifiers
// resulting from applying a filter to an IdentityProvider.
type IdentityFilterIdentifierProvider struct {
	filter           flow.IdentityFilter
	identityProvider module.IdentityProvider
}

func NewIdentityFilterIdentifierProvider(filter flow.IdentityFilter, identityProvider module.IdentityProvider) *IdentityFilterIdentifierProvider {
	return &IdentityFilterIdentifierProvider{filter, identityProvider}
}

func (p *IdentityFilterIdentifierProvider) Identifiers() flow.IdentifierList {
	return p.identityProvider.Identities(p.filter).NodeIDs()
}
