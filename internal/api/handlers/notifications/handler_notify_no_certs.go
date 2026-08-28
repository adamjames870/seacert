package notifications

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
)

func HandlerNotifyNoCertsEmail(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
