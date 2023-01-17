package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/io"
)

// LoadNodeMachineAccountInfoFile loads machine account info from the default location within the
// bootstrap directory - Currently being used by Collection and Consensus nodes
func LoadNodeMachineAccountInfoFile(bootstrapDir string, nodeID flow.Identifier) (*bootstrap.NodeMachineAccountInfo, error) {

	// attempt to read file
	machineAccountInfoPath := filepath.Join(bootstrapDir, fmt.Sprintf(bootstrap.PathNodeMachineAccountInfoPriv, nodeID))
	bz, err := io.ReadFile(machineAccountInfoPath)
	if err != nil {
		return nil, fmt.Errorf("could not read machine account info: %w", err)
	}

	// unmashal machine account info
	var machineAccountInfo bootstrap.NodeMachineAccountInfo
	err = json.Unmarshal(bz, &machineAccountInfo)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal machine account info: %w", err)
	}

	return &machineAccountInfo, nil
}
