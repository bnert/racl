package api

import (
    "fmt"
    "context"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    repo "github.com/bnert/racl/internal/racl_repo"
)

func getAcl(w http.ResponseWriter, r *http.Request) (int, data) {
    resourceId := r.URL.Query().Get("r")
    entityId := chi.URLParam(r, "entityId")
    if resourceId == "" {
        return 400, data{
            "error": "'r' url param missing. Please provide with resource param",
        }
    }

    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()

    acl, err := repo.New(conn).GetAclForResourceByEntity(
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
            "capabilities": acl.Capabilities,
        },
    }
}

type createAclBody struct {
    Resource     string `json:"resource"`
    Entity       string `json:"entity"`
    Capabilities *[]string `json:"capabilities, omitempty"`
}

func createAcl(w http.ResponseWriter, r *http.Request) (int, data) {
    var body createAclBody
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return 400, data{"error": err.Error()}
    }

    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()


    q := repo.New(conn)

    resource, err := q.CreateResource(
        context.Background(),
        body.Resource,
    )
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
        },
    }
}

type UpdateBody struct {
    Resource     string    `json:"resource"`
    Capabilities *[]string `json:"capabilities, omitempty"`
}

func updateAcl(w http.ResponseWriter, r *http.Request) (int, data) {
    var body UpdateBody 
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        fmt.Println(err.Error())
        return 400, data{"error": "Must provide 'capabilities' and 'resource'."}
    }

    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()

    entityId := chi.URLParam(r, "entityId")
    q := repo.New(conn)

    existingAcl, err := q.GetAclForResourceByEntity(
        context.Background(),
        repo.GetAclForResourceByEntityParams{
            ResourceID: body.Resource,
            Entity: entityId,
        },
    )
    if err != nil {
        return 404, data{"error": "entity not found"}
    }

    if body.Capabilities == nil {
        return 200, data{
            "data": data {
                "resource": existingAcl.ResourceID,
                "entity": existingAcl.Entity,
                "capabilities": existingAcl.Capabilities,
            },
        }
    }

    updateAcl, err := q.UpdateAclCapabilities(
        context.Background(),
        repo.UpdateAclCapabilitiesParams{
            Entity: existingAcl.Entity,
            Capabilities: *body.Capabilities,
        },
    )
    if err != nil {
        return 500, data{"error": err.Error()}
    }

    return 200, data{
        "data": data{
            "resource": updateAcl.ResourceID,
            "entity": updateAcl.Entity,
            "capabilities": updateAcl.Capabilities,
        },
        "meta": data{
            "capabilities": data{
                "prev": existingAcl.Capabilities, 
            },
        },
    }
}

type DeleteBody struct {
    Resource string `json:"resource"`
}

func deleteAcl(w http.ResponseWriter, r *http.Request) (int, data) {
    conn, err := ctxDbConn(r)
    if err != nil {
        return 500, data{"error": "unable to establish db connection"}
    }
    defer conn.Release()

    var body DeleteBody
    err = json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        fmt.Println(err.Error())
        return 400, data{"error": "Must provide 'resource'."}
    }

    entityId := chi.URLParam(r, "entityId")
    q := repo.New(conn)

    deletedAcl, err := q.DeleteAcl(
        context.Background(),
        repo.DeleteAclParams{
            Entity: entityId,
            ResourceID: body.Resource,
        },
    )
    if err != nil {
        return 404, data{"error": "Unable to find acl."}
    }

    return 200, data{
        "data": data{
            "resource": deletedAcl.ResourceID,
            "entity": deletedAcl.Entity,
        },
    }
}

