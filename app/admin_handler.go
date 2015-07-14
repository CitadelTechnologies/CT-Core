package app

import(
    "encoding/json"
    "gleipnir/errors"
    "net/http"
)

func sendStatus(w http.ResponseWriter, r *http.Request) {

    Core.refreshProfile()

    w.Header().Set("Access-Control-Allow-Origin", "*")
    err := json.NewEncoder(w).Encode(Core)
    errors.Check(err)


}

func shutdownKernel(w http.ResponseWriter, r *http.Request) {

    Core.Shutdown()
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)

}

func runKernel(w http.ResponseWriter, r *http.Request) {

    Core.Run()
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)

}

func shutdownServices(w http.ResponseWriter, r *http.Request) {

    Core.ShutdownServices()
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)

}