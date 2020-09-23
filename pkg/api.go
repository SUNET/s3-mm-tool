package api
import (
	"net/http"
	"crypto/tls"
	
	"github.com/gorilla/mux"
    "github.com/sirupsen/logrus"	
)

var Log = logrus.New()
var Handler *mux.Router

func Listen(hostPort string) {
        go func() {
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
        }()
}


func init() {
        Handler = mux.NewRouter()
        Handler.StrictSlash(true)

		Handler.HandleFunc("/api/create", create_dataset)
}

func create(w http.ResponseWriter, r *http.Request) {
	
}
