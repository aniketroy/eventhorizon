package writerhttp

import (
	"encoding/json"
	"github.com/function61/eventhorizon/writer"
	"github.com/function61/eventhorizon/writer/authmiddleware"
	wtypes "github.com/function61/eventhorizon/writer/types"
	"net/http"
)

func AppendToStreamHandlerInit(eventWriter *writer.EventstoreWriter) {
	ctx := eventWriter.GetConfigurationContext()

	http.Handle("/writer/append", authmiddleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appendToStreamRequest wtypes.AppendToStreamRequest
		if err := json.NewDecoder(r.Body).Decode(&appendToStreamRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		output, err := eventWriter.AppendToStream(appendToStreamRequest.Stream, appendToStreamRequest.Lines)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(output)
	}), ctx))
}
