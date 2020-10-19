package types

type Projects struct {
	Projects map[string]Project `json:"projects"`
}

type Project struct {
	Name string `json:"name"`
}
