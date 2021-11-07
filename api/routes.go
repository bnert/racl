package api

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    repo "github.com/bnert/racl/internal/racl_repo"
)

func getResource(w http.ResponseWriter, r *http.Request) (int, data) {
    // Expects the "r" query parameter, which is the resource

    resourceId := r.URL.Query().Get("r")
    entityId := chi.URLParam(r, "entityId")

    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()

    if resourceId == "" {
        return 400, data{
            "error": "'r' url param missing. Please provide with resource param",
        }
    }

    q := repo.New(conn)

    capabilities, err := q.GetAclForResourceByEntity(
        context.Background(),
        repo.GetAclForResourceByEntityParams{
            Entity: entityId,
            ResourceID: resourceId,
        },
    )
    if err != nil {
        return 404, data{"error": "Unable to find entity/resource acl"}
    }

    return 200, data{
        "data": data{
            "capabilities": capabilities,
        },
    }
}

type CreateResourceBody struct {
    Resource     string `json:"resource"`
    Entity       string `json:"entity"`
    Capabilities *[]string `json:"capabilities, omitempty"`
}

func createResource(w http.ResponseWriter, r *http.Request) (int, data) {
    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()

    var body CreateResourceBody
    err = json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return 400, data{"error": err.Error()}
    }

    q := repo.New(conn)

    resource, err := q.CreateResource(context.Background(), body.Resource)
    if err != nil {
        return 400, data{"error": err.Error()}
    }

    capabilities := []string{"c", "r", "u", "d", "a"}
    if body.Capabilities != nil {
        capabilities = *body.Capabilities
    }

    acl, err := q.CreateAcl(context.Background(), repo.CreateAclParams{
        ResourceID: resource.ID,
        Entity: body.Entity,
        Capabilities: capabilities,
    })

    return 201, data{
        "data": data{
            "resource": acl.ResourceID,
            "entity": acl.Entity,
            "capabilities": acl.Capabilities,
        }
    }
}

func updateResource(w http.ResponseWriter, r *http.Request) (int, data) {
    return 200, data{"data": "updated"}
}

func deleteResource(w http.ResponseWriter, r *http.Request) (int, data) {
    return 200, data{"data": "deleted"}
}
