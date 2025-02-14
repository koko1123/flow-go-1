package testutil

import (
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/koko1123/flow-go-1/fvm"
	"github.com/koko1123/flow-go-1/model/flow"
)

func CreateTokenTransferTransaction(chain flow.Chain, amount int, to flow.Address, signer flow.Address) *flow.TransactionBody {
	return flow.NewTransactionBody().
		SetScript([]byte(fmt.Sprintf(`
		import FungibleToken from 0x%s
		import FlowToken from 0x%s

		transaction(amount: UFix64, to: Address) {
			let sentVault: @FungibleToken.Vault

			prepare(signer: AuthAccount) {
				let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
					?? panic("Could not borrow reference to the owner's Vault!")
				self.sentVault <- vaultRef.withdraw(amount: amount)
			}

			execute {
				let receiverRef = getAccount(to)
					.getCapability(/public/flowTokenReceiver)
					.borrow<&{FungibleToken.Receiver}>()
					?? panic("Could not borrow receiver reference to the recipient's Vault")
				receiverRef.deposit(from: <-self.sentVault)
			}
		}`, fvm.FungibleTokenAddress(chain), fvm.FlowTokenAddress(chain)))).
		AddArgument(jsoncdc.MustEncode(cadence.UFix64(amount))).
		AddArgument(jsoncdc.MustEncode(cadence.NewAddress(to))).
		AddAuthorizer(signer)
}
