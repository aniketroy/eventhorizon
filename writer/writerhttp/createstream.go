package writerhttp

import (
	"encoding/json"
	"github.com/function61/pyramid/writer"
	"github.com/function61/pyramid/writer/authmiddleware"
	wtypes "github.com/function61/pyramid/writer/types"
	"io"
	"net/http"
)

func CreateStreamHandlerInit(eventWriter *writer.EventstoreWriter) {
	ctx := eventWriter.GetConfigurationContext()

	// $ curl -d '{"Name": "/foostream"}' http://localhost:9092/create_stream
	http.Handle("/create_stream", authmiddleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var createStreamRequest wtypes.CreateStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&createStreamRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := eventWriter.CreateStream(createStreamRequest.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "OK\n")
	}), ctx))
}
