# racl (`r`est `a`ccess `c`ontrol `l`ists)

## Motivation
Authorization can be hard, and this project aims to be simple solution to general
authz problems. Mainly, provide a standalone binary which can be configured in a
simple manner, deployed in a simple manner, and be simple to interface with. In other
words, I want to make authorization of generic resources simpler than what is currently
out there, and not have to hand roll authorzation per project.

## Tools
- Go
- sqlc

## High Level Design
Taking a note/cue from Hashicorp, the acl's available by racl will be implemented
around the notion of "capabilities", which are:

- c: create
- r: read
- u: update
- d: delete
- a: admin

Each of the capabilities listed above are a single character, and will be returned
as such. The 'a'/admin capability encapsulates all c/r/u/d operations and if is in
the acl will allow any operation. For instance, if an acl looks like: ['r', 'a'], then
the entity which is attached to said acl will have c/r/u/d capabilities on the
referenced resource.

### Authentication
One the first startup/connection to a database of the application, an api key and secret
will be created. They will be output to a file and be cryptographically strong. Therefore,
there will be no need to hash the api key/secret, which will make ther service more
performant than storing the secret as a hashed string.

### Operations
All operations shall be authenticated 

## API Design

### Create and ACL
By default, the entity which creates a resource will have c/r/u/d/a capabilities,
unless otherwise specified.

#### Create an acl with default capabilities
Request:
```
POST /acl/

{
  "resource": "some resource id",
  "entity": "<entity id>"
}
```

Response:
```
200 OK

{
  "data": {
    "id": "some uuid",
    "resource": "some resource id",
    "entity": "some entity id",
    "capabilities": ["c", "r", "u", "d", "a"]
  }
}
```

#### Create an acl and override default capabilities
Request:
```
POST /acl/

{
  "resource": "some resource id",
  "entity": "<entity id>",
  "capabilities: ["r"]
}
```

Response:
```
200 OK

{
  "data": {
    "id": "some uuid",
    "resource": "some resource id",
    "entity": "some entity id",
    "capabilities": ["r"]
  }
}
```

#### Updating (or creating) an acl

If you PUT as new resource, then the default will be the admin
roles (admin, create, read, update, delete), otherwise,
the capabilities provided will be respected. The resourceId __MUST__ be
a uuid v4, otherwise a 4XX response will be returned.

Request for resoure which exists:
```
PUT /acl/{entityId}

{
  "resource": "some resource id",
  "capabilities": ["c", "r", "u"]
}
```

Response for resource which exists:
```
200 OK

{
  "data": {
    "resource": "some resource id",
    "entity": "some entity id",
    "capabilities": ["c", "r", "u"]
  },
  "meta": {
    "capabilities": {
      "prev": ["r"]
    }
  }
}
```

Request for resource which does not exist:
```
PUT /acl/{entityId}

{
  "resource": "some resource id",
  "capabilities": ["c", "r", "u"]
}
```

Response for resource which does not exit:
```
201 Created

{
  "data": {
    "resource": "some resource id",
    "entity": "some entity id",
    "capabilities": ["c", "r", "u"]
  },
  "meta": {
    "capabilities": {
      "prev": []
    }
  }
}
```

### Query an ACL
Return the capabilities of an entity for a given resource
```
GET /acl/{entityId}?r={resourceId}

200
{
  "data": {
    "capabilities": ["c", "r", "u", "d"]
  }
}
```

Return if an operation is permitted for an entity on a resource
```
GET /acl/{entityId}?r={resourceId}&c=u


200
{
  "data": {
    "allowed": true
  }
}  
```

### Remove an acl for an entity
This will only respect the "delete" and "admin" capabilities.

Request to delete an entity attached to a resource:
```
DELETE /acl/{entityId}

{
  "resource": "some resource id",
}
```

Response to delete an entity attached to a resource:
```
200 OK

{
  "data": {
    "entity": "some entity id",
    "resource": "some resource id"
  }
}
```

### Remove a resource
This operation also removes all the reference entity acl's.

Request to delete a resource
```
DELETE /resource/{resourceId}
```

Response to delete a resource
```
200 OK

{
  "data": {
    "resource": "some resource id",
  }
}
```

