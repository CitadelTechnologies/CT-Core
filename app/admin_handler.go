package app

import(
    "encoding/json"
    "gleipnir/errors"
    "net/http"
)

func sendStatus(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "*")
    err := json.NewEncoder(w).Encode(Core)
    errors.Check(err)


}