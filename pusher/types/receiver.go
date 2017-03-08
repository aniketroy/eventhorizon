package types

import (
	rtypes "github.com/function61/eventhorizon/reader/types"
)

type Receiver interface {
	GetSubscriptionId() string
	PushReadResult(*rtypes.ReadResult) *PushResult
}
