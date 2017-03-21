package target

import (
	"encoding/json"
	"github.com/function61/pyramid/cli/example_target/target/schema"
	"net/http"
)

func (pa *Target) setupJsonRestApi() {
	http.Handle("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var users []schema.User
		err := pa.db.All(&users)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(users); err != nil {
			panic(err)
		}
	}))

	http.Handle("/companies", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var companies []schema.Company
		err := pa.db.All(&companies)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(companies); err != nil {
			panic(err)
		}
	}))
}
