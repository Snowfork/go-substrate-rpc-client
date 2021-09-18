// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"github.com/snowfork/go-substrate-rpc-client/v3/client"
	"github.com/snowfork/go-substrate-rpc-client/v3/types"
)

// GetKeys retreives the keys with the given prefix
func (s *State) GetKeys(prefix types.StorageKey, blockHash types.Hash) ([]types.StorageKey, error) {
	return s.getKeys(prefix, &blockHash)
}

// GetKeysLatest retreives the keys with the given prefix for the latest block height
func (s *State) GetKeysLatest(prefix types.StorageKey) ([]types.StorageKey, error) {
	return s.getKeys(prefix, nil)
}

// GetKeysPaged retreives the keys with the given prefix with paginated results
func (s *State) GetKeysPaged(prefix types.StorageKey, //nolint:interfacer
	count uint32, startKey *types.StorageKey, //nolint:interfacer
	blockHash types.Hash) ([]types.StorageKey, error) {
	var res []string

	var startKeyHex *string
	if startKey != nil {
		hex := startKey.Hex()
		startKeyHex = &hex
	}
	err := client.CallWithBlockHash(s.client, &res, "state_getKeysPaged", &blockHash, prefix.Hex(), count, startKeyHex)
	if err != nil {
		return nil, err
	}

	return decodeResponse(res)
}

func (s *State) getKeys(prefix types.StorageKey, blockHash *types.Hash) ([]types.StorageKey, error) {
	var res []string
	err := client.CallWithBlockHash(s.client, &res, "state_getKeys", blockHash, prefix.Hex())
	if err != nil {
		return nil, err
	}

	return decodeResponse(res)
}

func decodeResponse(response []string) ([]types.StorageKey, error) {
	keys := make([]types.StorageKey, len(response))
	for i, r := range response {
		err := types.DecodeFromHexString(r, &keys[i])
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
