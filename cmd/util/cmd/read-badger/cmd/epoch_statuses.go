package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/koko1123/flow-go-1/cmd/util/cmd/common"
	"github.com/koko1123/flow-go-1/model/flow"
)

func init() {
	rootCmd.AddCommand(epochStatusesCmd)

	epochStatusesCmd.Flags().StringVarP(&flagBlockID, "block-id", "b", "", "the block id of which to query the epoch status")
	_ = epochStatusesCmd.MarkFlagRequired("block-id")
}

var epochStatusesCmd = &cobra.Command{
	Use:   "epoch-statuses",
	Short: "get epoch statuses by block ID",
	Run: func(cmd *cobra.Command, args []string) {
		storages, db := InitStorages()
		defer db.Close()

		log.Info().Msgf("got flag block id: %s", flagBlockID)
		blockID, err := flow.HexStringToIdentifier(flagBlockID)
		if err != nil {
			log.Error().Err(err).Msg("malformed block id")
			return
		}

		log.Info().Msgf("getting epoch status by block id: %v", blockID)
		epochStatus, err := storages.Statuses.ByBlockID(blockID)
		if err != nil {
			log.Error().Err(err).Msgf("could not get epoch status for block id: %v", blockID)
			return
		}

		common.PrettyPrint(epochStatus)
	},
}
