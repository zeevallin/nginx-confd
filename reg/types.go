package main

type Mapping struct {
	Upstream string
	Pattern  string
	Port     int
}

type Server struct {
	Host   string            `json:"host"`
	Config map[string]string `json:"config,omitempty"`
}
