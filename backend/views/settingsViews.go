package views

import (
	"encoding/json"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/maf/interfaces"
	"net/http"
)

type SettingsView struct {
	SettingsService controllers.SettingsService `inject:""`
}

func (v *SettingsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	settings := v.SettingsService.Get()

	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(settings)

	if n != nil {
		n(rw, r, ctx)
	}

}

type SettingsUpdateView struct {
	SettingsService controllers.SettingsService `inject:""`
}

func (v *SettingsUpdateView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	settings := v.SettingsService.Get()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&settings)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	err = v.SettingsService.Update()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)

}
