package pushlib

// This small library receives pushes over HTTP + JSON and feeds them to your
// application via PushAdapter (you have to implement this), which essentially
// is a bridge between handling events and persisting the results to a database.

import (
	"encoding/json"
	"github.com/function61/pyramid/cursor"
	ptypes "github.com/function61/pyramid/pusher/types"
	"log"
	"net/http"
)

type Listener struct {
	subscriptionId string
	adapter        PushAdapter
}

func New(subscriptionId string, adapter PushAdapter) *Listener {
	return &Listener{
		subscriptionId: subscriptionId,
		adapter:        adapter,
	}
}

// called by the HTTP endpoint for pushing.
// returns PushOutput for sending status back to Pusher
func (l *Listener) Push(input *ptypes.PushInput) (*ptypes.PushOutput, error) {
	var output *ptypes.PushOutput

	// ask adapter to provide us with a transaction
	err := l.adapter.PushWrapTransaction(func(tx interface{}) error {
		var err error
		output, err = l.pushInternal(input, tx)

		return err
	})

	return output, err
}

func (l *Listener) pushInternal(input *ptypes.PushInput, tx interface{}) (*ptypes.PushOutput, error) {
	// ensure that subscription ID is correct
	if input.SubscriptionId != l.subscriptionId {
		return ptypes.NewPushOutputIncorrectSubscriptionId(l.subscriptionId), nil
	}

	fromOffset := cursor.CursorFromserializedMust(input.Read.FromOffset)

	// ensure that Pusher is continuing Push of the stream from the stream
	// offset that we last saved
	ourOffset := l.queryOffset(fromOffset.Stream, tx)

	if !fromOffset.PositionEquals(ourOffset) {
		return ptypes.NewPushOutputIncorrectBaseOffset(ourOffset.Serialize()), nil
	}

	// start with the offset stored in database. if we don't ACK a single
	// event, this is what we'll return and pusher will know that we didn't move
	// forward and throttle the pushes accordingly
	acceptedOffset := ourOffset.Serialize()

	behindCursors := make(map[string]string)

	for _, line := range input.Read.Lines {
		if line.IsMeta {
			// everything we encounter in SubscriptionActivity is something we ourselves
			// have subscribed to, so we can just check:
			// => if we're behind
			// => if we're never heard of the stream => start following it
			for _, remoteCursorSerialized := range line.SubscriptionActivity {
				remoteCursor := cursor.CursorFromserializedMust(remoteCursorSerialized)

				// see if this stream's behind-ness is already confirmed as behind?
				// in that case we don't need newer data because we already know our
				// position for this stream, and re-checking it will never change it.
				_, weAlreadyKnowThisStreamIsBehind := behindCursors[remoteCursor.Stream]

				if !weAlreadyKnowThisStreamIsBehind {
					shouldStartFrom := l.isRemoteAhead(remoteCursor, tx)

					if shouldStartFrom != nil {
						log.Printf("Listener: remote ahead of us: %s", remoteCursorSerialized)

						behindCursors[remoteCursor.Stream] = shouldStartFrom.Serialize()
					}
				}
			}
		} else {
			if line.Content == "" {
				log.Printf("Pusher: empty line %v", line)
			}

			if err := l.adapter.PushHandleEvent(line.Content, tx); err != nil {
				return nil, err
			}
		}

		// only ACK offsets if no behind streams encountered
		// (this happens only for subscription streams anyway)
		if len(behindCursors) == 0 {
			acceptedOffset = line.PtrAfter
		}
	}

	if err := l.adapter.PushSetOffset(fromOffset.Stream, acceptedOffset, tx); err != nil {
		return nil, err
	}

	return ptypes.NewPushOutputSuccess(acceptedOffset, stringMapToSlice(behindCursors)), nil
}

func (l *Listener) queryOffset(stream string, tx interface{}) *cursor.Cursor {
	cursorSerialized, exists := l.adapter.PushGetOffset(stream, tx)

	// we can trust that it is a valid stream because all pushes are based on
	// the subscription ID that is exclusive to us. so if stream does not exist
	// => allow it to be created
	if !exists {
		return cursor.BeginningOfStream(stream, cursor.UnknownServer)
	}

	return cursor.CursorFromserializedMust(cursorSerialized)
}

func (l *Listener) isRemoteAhead(remote *cursor.Cursor, tx interface{}) *cursor.Cursor {
	ourCursor := l.queryOffset(remote.Stream, tx)

	if remote.IsAheadComparedTo(ourCursor) {
		return ourCursor
	} else {
		return nil
	}
}

func (l *Listener) AttachPushHandler() {
	http.Handle("/_pyramid_push", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var push ptypes.PushInput
		if err := json.NewDecoder(r.Body).Decode(&push); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		output, err := l.Push(&push)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(output)
	}))
}

func stringMapToSlice(map_ map[string]string) []string {
	slice := []string{}

	for _, value := range map_ {
		slice = append(slice, value)
	}

	return slice
}
