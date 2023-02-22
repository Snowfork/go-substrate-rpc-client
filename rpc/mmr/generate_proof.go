package mmr

import (
	"github.com/snowfork/go-substrate-rpc-client/v4/client"
	"github.com/snowfork/go-substrate-rpc-client/v4/types"
)

// GenerateProof retrieves a MMR proof and leaf for the specified leave index, at the given blockHash (useful to query a
// proof at an earlier block, likely with antoher MMR root)
func (c *MMR) GenerateProof(blockNumber uint32, blockHash types.Hash) (types.GenerateMMRProofResponse, error) {
	return c.generateProof(blockNumber, &blockHash)
}

// GenerateProofLatest retrieves the latest MMR proof and leaf for the specified leave index
func (c *MMR) GenerateProofLatest(blockNumber uint32) (types.GenerateMMRProofResponse, error) {
	return c.generateProof(blockNumber, nil)
}

func (c *MMR) generateProof(blockNumber uint32, blockHash *types.Hash) (types.GenerateMMRProofResponse, error) {
	var res types.GenerateMMRProofResponse
	blocks := [1]uint32{blockNumber}
	err := client.CallWithBlockHash(c.client, &res, "mmr_generateProof", blockHash, blocks, blockNumber)
	if err != nil {
		return types.GenerateMMRProofResponse{}, err
	}

	return res, nil
}
