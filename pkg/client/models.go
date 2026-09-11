package client

type Login struct {
	Login        string   `json:"login"`
	Name         string   `json:"name"`
	Password     string   `json:"password"`
	UUID         string   `json:"uuid"`
	StringFields []string `json:"string_fields"`
	Group        string   `json:"group"`
	Totp         string   `json:"totp,omitempty"`
}

type AssociateResponse struct {
	AssociatedName    string `json:"associated_name"`
	IdentificationKey string `json:"identification_key"`
}
