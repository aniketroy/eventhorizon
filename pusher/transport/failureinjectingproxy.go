package transport

import (
	"errors"
	"fmt"
	ptypes "github.com/function61/eventhorizon/pusher/types"
)

// this tranport proxies to another transport ("endpoint"), but injects failures
// to some the pushes. This is used to test resiliency against failures.

type FailureInjectingProxy struct {
	endpoint ptypes.Transport
	counter  int
}

func NewFailureInjectingProxy(endpoint ptypes.Transport) *FailureInjectingProxy {
	return &FailureInjectingProxy{
		endpoint: endpoint,
		counter:  1, // so first req fails
	}
}

func (f *FailureInjectingProxy) Push(input *ptypes.PushInput) (*ptypes.PushOutput, error) {
	if f.shouldFail() {
		return nil, errors.New(fmt.Sprintf("synthetic failure %d", f.counter))
	}

	return f.endpoint.Push(input)
}

func (f *FailureInjectingProxy) shouldFail() bool {
	defer func() { f.counter++ }()

	if f.counter%4 == 0 {
		return false // make every 4th request succeed
	}

	return true
}
