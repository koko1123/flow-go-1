package testnet

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"

	"github.com/koko1123/flow-go-1/integration/client"
	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/state/protocol/inmem"
	"github.com/koko1123/flow-go-1/utils/io"
	"github.com/onflow/flow-go/crypto"
)

// healthcheckAccessGRPC returns a Docker healthcheck function that pings the Access node GRPC
// service exposed at the given port.
func healthcheckAccessGRPC(apiPort string) func() error {
	return func() error {
		fmt.Println("healthchecking...")
		c, err := client.NewAccessClient(fmt.Sprintf(":%s", apiPort))
		if err != nil {
			return err
		}

		return c.Ping(context.Background())
	}
}

// healthcheckExecutionGRPC returns a Docker healthcheck function that pings the Execution node GRPC
// service exposed at the given port.
func healthcheckExecutionGRPC(apiPort string) func() error {
	return func() error {
		fmt.Println("healthchecking...")
		c, err := client.NewExecutionClient(fmt.Sprintf(":%s", apiPort))
		if err != nil {
			return err
		}

		return c.Ping(context.Background())
	}
}

// currentUser returns a uid:gid Unix user identifier string for the current
// user. This is used to run node containers under the same user to avoid
// permission conflicts on files mounted from the host.
func currentUser() string {
	cur, _ := user.Current()
	return fmt.Sprintf("%s:%s", cur.Uid, cur.Gid)
}

func toParticipants(confs []ContainerConfig) flow.IdentityList {
	il := make(flow.IdentityList, 0, len(confs))
	for _, conf := range confs {
		il = append(il, conf.Identity())
	}
	return il
}

func toNodeInfos(confs []ContainerConfig) []bootstrap.NodeInfo {
	infos := make([]bootstrap.NodeInfo, 0, len(confs))
	for _, conf := range confs {
		infos = append(infos, conf.NodeInfo)
	}
	return infos
}

func getSeed() ([]byte, error) {
	seedLen := int(math.Max(crypto.SeedMinLenDKG, crypto.KeyGenSeedMinLenBLSBLS12381))
	seed := make([]byte, seedLen)
	n, err := rand.Read(seed)
	if err != nil || n != seedLen {
		return nil, err
	}
	return seed, nil
}

func WriteJSON(path string, data interface{}) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	marshaled, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return WriteFile(path, marshaled)
}

func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	return err
}

// rootProtocolJsonWithoutAddresses strips out all node addresses from the root protocol json file specified as srcFile
// and creates the dstFile with the modified contents
func rootProtocolJsonWithoutAddresses(srcfile string, dstFile string) error {

	data, err := io.ReadFile(filepath.Join(srcfile))
	if err != nil {
		return err
	}

	var rootSnapshot inmem.EncodableSnapshot
	err = json.Unmarshal(data, &rootSnapshot)
	if err != nil {
		return err
	}

	strippedSnapshot := inmem.StrippedInmemSnapshot(rootSnapshot)

	return WriteJSON(dstFile, strippedSnapshot)
}
