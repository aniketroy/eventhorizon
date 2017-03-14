package metaevents

import (
	"github.com/function61/pyramid/util/ass"
	"testing"
)

func TestUnsubscribed(t *testing.T) {
	isMeta, _, event := Parse(".Unsubscribed {\"subscription_id\":\"6894605c-2a8e\",\"ts\":\"2017-02-27T17:12:31.446Z\"}")

	ass.True(t, isMeta)

	unsubscribed := event.(Unsubscribed)

	ass.EqualString(t, unsubscribed.SubscriptionId, "6894605c-2a8e")
}
