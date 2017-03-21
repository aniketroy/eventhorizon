package target

import (
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/function61/pyramid/pusher/pushlib"
	"log"
	"net/http"
)

// implements PushAdapter interface
type Target struct {
	pushListener *pushlib.Listener
	db           *storm.DB
	tx           *bolt.Tx
}

func NewTarget() *Target {
	subscriptionId := "foo"

	db, err := storm.Open("/tmp/listener.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	pa := &Target{
		db: db,
	}

	pa.pushListener = pushlib.New(
		subscriptionId,
		pa)

	return pa
}

func (pa *Target) Run() {
	pa.setupJsonRestApi()

	go pushlib.StartChildProcess("http://127.0.0.1:8080/_pyramid_push")

	// sets up HTTP endpoint for receiving pushes
	pa.pushListener.AttachPushHandler()

	srv := &http.Server{Addr: ":8080"}

	log.Printf("Target: listening at :8080")

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("Target: ListenAndServe() error: %s", err)
	}
}

func (pa *Target) PushGetOffset(stream string) (string, bool) {
	offset := ""
	if err := pa.db.WithTransaction(pa.tx).Get("cursors", stream, &offset); err != nil {
		if err == storm.ErrNotFound {
			return "", false
		}

		// more serious error
		panic(err)
	}

	return offset, true
}

func (pa *Target) PushSetOffset(stream string, offset string) {
	if err := pa.db.WithTransaction(pa.tx).Set("cursors", stream, offset); err != nil {
		panic(err)
	}
}

// this is where all the magic happens. pushlib calls this function for every
// incoming event from Pyramid.
func (pa *Target) PushHandleEvent(eventSerialized string) error {
	handled, err := applySerializedEvent(eventSerialized, pa)

	if err != nil {
		return err
	}

	if !handled {
		log.Printf("Target: unknown event: %s", eventSerialized)
	}

	return nil
}

func (pa *Target) PushTransaction(run func() error) error {
	// PushTransaction() is an API that pushlib calls to wrap all the following
	// operations in a single transaction. we:
	//
	//     1) start transaction
	//     2) store it in our own state (we use it from below mentioned APIs)
	//     3) call back to pushlib with "run" which will start rapidly calling
	//        PushHandleEvent() multiple times + PushSetOffset() once
	//     4) we get back error state from "run" callback indicating if anything went
	//        wrong. if we get error Bolt rollbacks the TX, if all is fine we commit.
	err := pa.db.Bolt.Update(func(tx *bolt.Tx) error {
		pa.tx = tx

		return run()
	})

	return err
}
