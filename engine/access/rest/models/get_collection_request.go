package models

import (
	"github.com/onflow/flow-go/engine/access/rest"
	"github.com/onflow/flow-go/model/flow"
)

const id = "id"
const expandsTransactions = "transactions"

type GetCollectionRequest struct {
	ID                  flow.Identifier
	ExpandsTransactions bool
}

func (g *GetCollectionRequest) Build(r *rest.Request) error {
	err := g.Parse(
		r.GetVar("id"),
	)
	if err != nil {
		return rest.NewBadRequestError(err)
	}

	g.ExpandsTransactions = r.Expands(expandsTransactions)

	return nil
}

func (g *GetCollectionRequest) Parse(rawID string) error {
	var id ID
	err := id.Parse(rawID)
	if err != nil {
		return err
	}

	g.ID = id.Flow()
	return nil
}
