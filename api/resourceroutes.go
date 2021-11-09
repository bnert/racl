package api

import (
    "net/http"
)

func getResource(w http.ResponseWriter, r *http.Request) (int, data) {
    return 200, data{"data": "not implemented"}
}

func createResource(w http.ResponseWriter, r *http.Request) (int, data) {
    return 200, data{"data": "not implemented"}
}

func deleteResource(w http.ResponseWriter, r *http.Request) (int, data) {
    return 200, data{"data": "not implemented"}
}
