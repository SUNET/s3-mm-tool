package api

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "github.com/sunet/s3-mm-tool/docs"
	"github.com/sunet/s3-mm-tool/pkg/index"
	"github.com/sunet/s3-mm-tool/pkg/manifest"
	"github.com/sunet/s3-mm-tool/pkg/meta"
	"github.com/sunet/s3-mm-tool/pkg/storage"
	"github.com/sunet/s3-mm-tool/pkg/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title S3 Metadata Manager Tool API
// @version 1.0
// @description Create and initialize a research data bucket with a provided metadata manifest configuration.
// @description Set ACL to support remote indexing and register dataset with index server.

// @contact.name SUNET NOC
// @contact.url https://www.sunet.se/
// @contact.email noc@sunet.se

// @license.name BSD

var Log = logrus.New()
var Handler *mux.Router

// CreateRequest example
type CreateRequest struct {
	Manifest manifest.ManifestInfo `json:"manifest"`
	Name     string                `json:"name" example:"mybucket"`
	Storage  storage.Destination   `json:"storage"`
}

// CreateResponse example
type CreateResponse struct {
	Manifest manifest.ManifestInfo `json:"manifest"`
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

// StatusResponse example
type StatusResponse struct {
	Name    string `json:"name" example:"s3-mm-tool"`
	Version string `json:"version" example:"1.0"`
}

// Status godoc
// @Summary Display status and version information
// @Description Display status and version information
// @Tags status
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get]
func Status(w http.ResponseWriter, r *http.Request) {
	status := StatusResponse{
		Name:    "s3-mm-tool",
		Version: meta.Version(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// CreateDataset godoc
// @Summary Creaate a new dataset bucket
// @Description Create a new bucket using supplied credentials and initialize with manifest
// @Tags create
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Create dataset"
// @Success 200 {object} CreateResponse
// @Failure 400
// @Failure 500
// @Router /create [post]
func CreateDataset(w http.ResponseWriter, r *http.Request) {
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

	err = index.RegisterURL(req.Storage.Endpoint + req.Name + "/.metadata/manifest.jsonld")

	var res CreateResponse
	res.Manifest = *m

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func init() {
	Handler = mux.NewRouter()
	Handler.StrictSlash(false)

	baseURL := os.Getenv("S3_MM_BASEURL")
	if len(baseURL) == 0 {
		baseURL = "http://localhost:3000"
	}

	Handler.HandleFunc("/create", CreateDataset)
	Handler.HandleFunc("/status", Status)
	Handler.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(baseURL+"/swagger/doc.json"),
		httpSwagger.DeepLinking(false),
	))
}
