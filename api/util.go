package api

import (
    "encoding/json"
    "net/http"
)

type data = map[string]interface{}

func handler(h func(w http.ResponseWriter, r *http.Request) (int, data)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        status, data := h(w, r)
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(data)
    }
}
