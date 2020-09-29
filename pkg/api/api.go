package api

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/sunet/s3-mm-tool/pkg/manifest"
	"github.com/sunet/s3-mm-tool/pkg/meta"
	"github.com/sunet/s3-mm-tool/pkg/utils"
)

var Log = logrus.New()
var Handler *mux.Router

func Listen(hostPort string) {
	http.Handle("/", Handler)
	cfg := utils.GetTLSConfig()
	if cfg == nil {
		Log.Fatal(http.ListenAndServe(hostPort, nil))
	} else {
		Log.Debug("Enabling TLS...")
		listener, err := tls.Listen("tcp", hostPort, cfg)
		if err != nil {
			Log.Fatal(err)
		}
		Log.Fatal(http.Serve(listener, nil))
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"name":    "s3-mm-tool",
		"version": meta.Version(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func create_dataset(w http.ResponseWriter, r *http.Request) {
	m, err := manifest.NewManifest()
	var mreq manifest.Manifest

	err = json.NewDecoder(r.Body).Decode(&mreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		manifest.MergeManifests(m, &mreq)
	}
	Log.Debug(*m)
}

func init() {
	Handler = mux.NewRouter()
	Handler.StrictSlash(true)

	Handler.HandleFunc("/api/create", create_dataset)
	Handler.HandleFunc("/api/status", status)
}
