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

package types

import "github.com/snowfork/go-substrate-rpc-client/v4/scale"

type SignedBlock struct {
	Block          Block                `json:"block"`
	Justifications OptionJustifications `json:"justifications"`
}

// Block encoded with header and extrinsics
type Block struct {
	Header     Header
	Extrinsics []Extrinsic
}

type Justification struct {
	ConsensusEngineID [4]byte
	EncodedJustification []byte
}

type OptionJustifications struct {
	option
	value []Justification
}

func (o OptionJustifications) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionJustifications) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionJustifications) SetSome(value []Justification) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionJustifications) SetNone() {
	o.hasValue = false
	o.value = []Justification{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionJustifications) Unwrap() (ok bool, value []Justification) {
	return o.hasValue, o.value
}
