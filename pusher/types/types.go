package types

import (
	rtypes "github.com/function61/pyramid/reader/types"
)

const (
	CodeSuccess                 = "success"
	CodeIncorrectBaseOffset     = "incorrect_base_offset"
	CodeIncorrectSubscriptionId = "incorrect_subscription_id"
)

type PushInput struct {
	SubscriptionId string
	// TODO: rename to push, fix in wire protocol doc
	// also link to this file from it
	Read           *rtypes.ReadResult
}

func NewPushInput(subscriptionId string, readResult *rtypes.ReadResult) *PushInput {
	return &PushInput{
		SubscriptionId: subscriptionId,
		Read:           readResult,
	}
}

type PushOutput struct {
	Code                  string
	AcceptedOffset        string
	CorrectSubscriptionId string // omit if empty?
	BehindCursors         []string
}

func NewPushOutputIncorrectBaseOffset(correctBaseOffset string) *PushOutput {
	return &PushOutput{
		Code:           CodeIncorrectBaseOffset,
		AcceptedOffset: correctBaseOffset,
	}
}

func NewPushOutputIncorrectSubscriptionId(correctSubscriptionId string) *PushOutput {
	return &PushOutput{
		Code: CodeIncorrectSubscriptionId,
		CorrectSubscriptionId: correctSubscriptionId,
	}
}

func NewPushOutputSuccess(acceptedOffset string, behindCursors []string) *PushOutput {
	return &PushOutput{
		Code:           CodeSuccess,
		AcceptedOffset: acceptedOffset,
		BehindCursors:  behindCursors,
	}
}
