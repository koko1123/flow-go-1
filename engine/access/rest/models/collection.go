package models

import (
	"fmt"

	"github.com/koko1123/flow-go-1/engine/access/rest/util"
	"github.com/koko1123/flow-go-1/model/flow"
)

const ExpandsTransactions = "transactions"

func (c *Collection) Build(
	collection *flow.LightCollection,
	txs []*flow.TransactionBody,
	link LinkGenerator,
	expand map[string]bool) error {

	self, err := SelfLink(collection.ID(), link.CollectionLink)
	if err != nil {
		return err
	}

	var expandable CollectionExpandable
	var transactions Transactions
	if expand[ExpandsTransactions] {
		transactions.Build(txs, link)
	} else {
		expandable.Transactions = make([]string, len(collection.Transactions))
		for i, tx := range collection.Transactions {
			expandable.Transactions[i], err = link.TransactionLink(tx)
			if err != nil {
				return err
			}
		}
	}

	c.Id = collection.ID().String()
	c.Transactions = transactions
	c.Links = self
	c.Expandable = &expandable

	return nil
}

func (c *CollectionGuarantee) Build(guarantee *flow.CollectionGuarantee) {
	c.CollectionId = guarantee.CollectionID.String()
	c.SignerIndices = fmt.Sprintf("%x", guarantee.SignerIndices)
	c.Signature = util.ToBase64(guarantee.Signature.Bytes())
}

type CollectionGuarantees []CollectionGuarantee

func (c *CollectionGuarantees) Build(guarantees []*flow.CollectionGuarantee) {
	collGuarantees := make([]CollectionGuarantee, len(guarantees))
	for i, g := range guarantees {
		var col CollectionGuarantee
		col.Build(g)
		collGuarantees[i] = col
	}
	*c = collGuarantees
}
