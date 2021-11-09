package api

import (
    "fmt"
    "context"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v4/pgxpool"
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

    acl, err := repo.New(conn).GetAclByEntityAndResource(
        context.Background(),
        repo.GetAclByEntityAndResourceParams{
            Entity: entityId,
            ResourceID: resourceId,
        },
    )
    if err != nil {
        return 404, data{"error": "Unable to find entity/resource acl."}
    }

    return 200, data{
        "data": data{
            "capabilities": acl.Capabilities,
        },
    }
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
    if err != nil {
        return 202, data{
            "data": data{
                "resourceID": acl.ResourceID,
                "entity": nil,
                "capabilities": []string{},
            },
            "error": "Unable to create acl. Retry operation.",
        }
    }

    return 201, data{
        "data": data{
            "resource": acl.ResourceID,
            "entity": acl.Entity,
            "capabilities": acl.Capabilities,
        },
    }
}

func createResourceAndAttachAcl(conn *pgxpool.Conn, entityId string, body *updateAclBody) (int, data) {
    q := repo.New(conn)

    _, err := q.CreateResource(
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

    acl, err := q.CreateAcl(
        context.Background(),
        repo.CreateAclParams{
            ResourceID: body.Resource,
            Entity: entityId,
            Capabilities: capabilities,
        },
    )
    if err != nil {
        return 202, data{
            "data": data{
                "resourceID": acl.ResourceID,
                "entity": nil,
                "capabilities": []string{},
            },
            "error": "Unable to create acl. Retry operation.",
        }
    }

    return 201, data{
        "data": data{
            "resourceID": acl.ResourceID,
            "entity": acl.Entity,
            "capabilities": acl.Capabilities,
        },
    }
}

func updateAcl(w http.ResponseWriter, r *http.Request) (int, data) {
    var body updateAclBody
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

    existingAcl, err := q.GetAclByEntityAndResource(
        context.Background(),
        repo.GetAclByEntityAndResourceParams{
            ResourceID: body.Resource,
            Entity: entityId,
        },
    )
    if err != nil {
        return createResourceAndAttachAcl(conn, entityId, &body)
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

    fmt.Println("CAPA", *body.Capabilities)

    updateAcl, err := q.UpdateAclCapabilities(
        context.Background(),
        repo.UpdateAclCapabilitiesParams{
            Entity: existingAcl.Entity,
            Capabilities: *body.Capabilities,
            ResourceID: existingAcl.ResourceID,
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

