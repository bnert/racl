package api

type createAclBody struct {
    Resource     string `json:"resource"`
    Entity       string `json:"entity"`
    Capabilities *[]string `json:"capabilities, omitempty"`
}

type updateAclBody struct {
    Resource     string    `json:"resource"`
    Capabilities *[]string `json:"capabilities, omitempty"`
}

