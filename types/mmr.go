package types

import (
	"encoding/json"
	"fmt"
)

// GenerateMMRProofResponse contains the generate proof rpc response
type GenerateMMRProofResponse struct {
	BlockHash H256
	Leaf      MMRLeaf
	Proof     MMRProof
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *GenerateMMRProofResponse) UnmarshalJSON(bz []byte) error {
	var tmp struct {
		BlockHash string `json:"blockHash"`
		Leaves    string `json:"leaves"`
		Proof     string `json:"proof"`
	}
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}
	err := DecodeFromHexString(tmp.BlockHash, &d.BlockHash)
	if err != nil {
		return err
	}
	var encodedLeaf []MMREncodableOpaqueLeaf
	err = DecodeFromHexString(tmp.Leaves, &encodedLeaf)
	if err != nil {
		return err
	}
	if len(encodedLeaf) == 0 {
		return fmt.Errorf("decode leaf error")
	}

	err = DecodeFromBytes(encodedLeaf[0], &d.Leaf)
	if err != nil {
		return err
	}
	var proof MultiMMRProof
	err = DecodeFromHexString(tmp.Proof, &proof)
	if err != nil {
		return err
	}
	if proof.LeafIndices == nil || len(proof.LeafIndices) == 0 {
		return fmt.Errorf("decode proof LeafIndices error")
	}
	d.Proof.LeafCount = proof.LeafCount
	d.Proof.Items = proof.Items
	d.Proof.LeafIndex = proof.LeafIndices[0]
	return nil
}

type ProofItem struct {
	Position U64
	Hash     H256
}

type ProofItemMarshal struct {
	Position uint64
	Hash     string
}

type GenerateAncestryProofResponse struct {
	PrevPeaks     []H256
	PrevLeafCount U64
	LeafCount     U64
	Items         []ProofItem
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *GenerateAncestryProofResponse) UnmarshalJSON(bz []byte) error {
	var tmp struct {
		PrevPeaks     []string         `json:"prev_peaks"`
		PrevLeafCount uint64           `json:"prev_leaf_count"`
		LeafCount     uint64           `json:"leaf_count"`
		Items         [][2]interface{} `json:"items"`
	}
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return fmt.Errorf("unmarshal JSON: %w", err)
	}

	d.PrevPeaks = make([]H256, len(tmp.PrevPeaks))
	for i, prevPeak := range tmp.PrevPeaks {
		err := DecodeFromHexString(prevPeak, &d.PrevPeaks[i])
		if err != nil {
			return err
		}
	}

	d.PrevLeafCount = NewU64(tmp.PrevLeafCount)
	d.LeafCount = NewU64(tmp.LeafCount)

	d.Items = make([]ProofItem, len(tmp.Items))
	for i, item := range tmp.Items {
		if len(item) != 2 {
			return fmt.Errorf("invalid item %d: expected [position, hash], got %v", i, item)
		}

		// Extract position (JSON number unmarshals as float64)
		position, ok := item[0].(float64)
		if !ok {
			return fmt.Errorf("invalid position in item %d: expected number, got %v", i, item[0])
		}

		// Extract hash (string)
		hash, ok := item[1].(string)
		if !ok {
			return fmt.Errorf("invalid hash in item %d: expected string, got %v", i, item[1])
		}

		// Assign to d.Items
		d.Items[i].Position = NewU64(uint64(position))
		if err := DecodeFromHexString(hash, &d.Items[i].Hash); err != nil {
			return fmt.Errorf("decode hash in item %d: %w", i, err)
		}
	}
	return nil
}

type MMREncodableOpaqueLeaf Bytes

// MMRProof is a MMR proof
type MMRProof struct {
	// The index of the leaf the proof is for.
	LeafIndex U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

// MultiMMRProof
type MultiMMRProof struct {
	// The indices of leaves.
	LeafIndices []U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

type MMRLeaf struct {
	Version               MMRLeafVersion
	ParentNumberAndHash   ParentNumberAndHash
	BeefyNextAuthoritySet BeefyNextAuthoritySet
	ParachainHeads        H256
}

type MMRLeafVersion U8

type ParentNumberAndHash struct {
	ParentNumber U32
	Hash         Hash
}

type BeefyNextAuthoritySet struct {
	// ID
	ID U64
	// Number of validators in the set.
	Len U32
	// Merkle Root Hash build from BEEFY uncompressed AuthorityIds.
	Root H256
}
