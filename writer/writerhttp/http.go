package writerhttp

import (
	"encoding/json"
	"github.com/function61/eventhorizon/config"
	"github.com/function61/eventhorizon/cursor"
	"github.com/function61/eventhorizon/reader"
	"github.com/function61/eventhorizon/writer"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ReadRequest struct {
	Cursor string
}

type CreateStreamRequest struct {
	Name string
}

type SubscribeToStreamRequest struct {
	Stream         string
	SubscriptionId string
}

type UnsubscribeFromStreamRequest struct {
	Stream         string
	SubscriptionId string
}

func HttpServe(eventWriter *writer.EventstoreWriter, shutdown chan bool, done chan bool) {
	reader := reader.NewEventstoreReader()

	srv := &http.Server{Addr: ":" + strconv.Itoa(config.WRITER_HTTP_PORT)}

	// $ curl -d '{"Cursor": "/tenants/foo:0:0"} http://localhost:9092/read
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		var readRequest ReadRequest
		if err := json.NewDecoder(r.Body).Decode(&readRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cur := cursor.CursorFromserializedMust(readRequest.Cursor)
		readResult, err := reader.Read(cur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		encoder.Encode(readResult)
	})

	// $ curl -d '{"Name": "/foostream"}' http://localhost:9092/create_stream
	http.HandleFunc("/create_stream", func(w http.ResponseWriter, r *http.Request) {
		var createStreamRequest CreateStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&createStreamRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := eventWriter.CreateStream(createStreamRequest.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "OK\n")
	})

	// $ curl -d '{"Stream": "/foostream", "SubscriptionId": "88c20701"}' http://localhost:9092/subscribe
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		var subscribeToStreamRequest SubscribeToStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&subscribeToStreamRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := eventWriter.SubscribeToStream(subscribeToStreamRequest.Stream, subscribeToStreamRequest.SubscriptionId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "OK\n")
	})

	// $ curl -d '{"Stream": "/foostream", "SubscriptionId": "88c20701"}' http://localhost:9092/unsubscribe
	http.HandleFunc("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		var unsubscribeFromStreamRequest UnsubscribeFromStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&unsubscribeFromStreamRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := eventWriter.UnsubscribeFromStream(unsubscribeFromStreamRequest.Stream, unsubscribeFromStreamRequest.SubscriptionId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "OK\n")
	})

	go func() {
		log.Printf("WriterHttp: binding to %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("WriterHttp: ListenAndServe() error: %s", err)
		}
	}()

	go func() {
		<-shutdown

		log.Printf("WriterHttp: shutting down")

		if err := srv.Shutdown(nil); err != nil {
			panic(err) // failed shutting down
		}

		log.Printf("WriterHttp: shutting down done")

		done <- true
	}()
}
