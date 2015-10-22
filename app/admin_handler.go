package app

import(
    "encoding/json"
    "net/http"
)

func sendStatus(w http.ResponseWriter, r *http.Request) {

    Core.refreshProfile()

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    err := json.NewEncoder(w).Encode(Core)
    CheckError(err)
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

    Core.ShutdownServices(false)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)

}