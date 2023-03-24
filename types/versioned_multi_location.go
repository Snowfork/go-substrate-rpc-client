package types

import (
	"fmt"

	"github.com/snowfork/go-substrate-rpc-client/v4/scale"
)

type VersionedMultiLocation struct {
	IsV3            bool
	MultiLocationV3 MultiLocationV3
}

func (m *VersionedMultiLocation) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 3:
		m.IsV3 = true
		return decoder.Decode(&m.MultiLocationV3)
	}

	return fmt.Errorf("unsupported variant: %d", b)
}

func (m VersionedMultiLocation) Encode(encoder scale.Encoder) error {
	switch {
	case m.IsV3:
		if err := encoder.PushByte(3); err != nil {
			return err
		}
		return encoder.Encode(m.MultiLocationV3)
	}

	return fmt.Errorf("unsupported variant")
}
