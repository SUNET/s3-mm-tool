package api

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/sunet/s3-mm-tool/pkg/manifest"
	"github.com/sunet/s3-mm-tool/pkg/meta"
	"github.com/sunet/s3-mm-tool/pkg/storage"
	"github.com/sunet/s3-mm-tool/pkg/utils"
)

var Log = logrus.New()
var Handler *mux.Router

type CreateRequest struct {
	Manifest manifest.Manifest   `json:"manifest"`
	Name     string              `json:"name"`
	Storage  storage.Destination `json:"storage"`
}

type CreateResponse struct {
	Manifest manifest.Manifest `json:"manifest"`
}

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
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req CreateRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.MergeManifests(req.Manifest)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.CreateS3DataContainer(req.Storage, req.Name, *m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res CreateResponse
	res.Manifest = *m

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func init() {
	Handler = mux.NewRouter()
	Handler.StrictSlash(true)

	Handler.HandleFunc("/api/create", create_dataset)
	Handler.HandleFunc("/api/status", status)
}
