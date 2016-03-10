package views

import (
	"encoding/json"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/maf/interfaces"
	"net/http"
	"time"
)

const dateLayout = "02-01-2006"

type StatsView struct {
	StatsController controllers.StatsController `inject:""`
}

func (v *StatsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse GET values
	fromRaw := r.FormValue("from")
	toRaw := r.FormValue("to")

	from, err := time.Parse(dateLayout, fromRaw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	to, err := time.Parse(dateLayout, toRaw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Loads stats object
	stats, err := v.StatsController.LoadStats(from, to)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Return stats object as JSON
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	enc.Encode(stats)

	if n != nil {
		n(rw, r, ctx)
	}

}
