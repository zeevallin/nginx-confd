package main

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"os"
	"time"
)

type Site struct {
	Listen     string `json:"listen"`
	ServerName string `json:"server_name"`

	Locations map[string]Location `json:"locations,omitempty"`

	Includes []string          `json:"includes,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

type Location struct {
	Upstream string `json:"upstream"`

	Includes []string          `json:"includes,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

func NewSite(name string) Site {
	return Site{
		Listen:     "80",
		ServerName: name,
		Locations: map[string]Location{
			"/": Location{
				Upstream: "app",
				Includes: []string{
					"default/proxy-headers",
				},
			},
		},
	}
}

func FormatSite(site Site) string {
	marshal, _ := json.Marshal(site)
	value := string(marshal)
	return value
}

func AnnounceSite(id string, site Site, client *etcd.Client) {
	key := fmt.Sprintf("/nginx/sites/%s", id)
	value := FormatSite(site)

	client.Set(key, value, 60)
	log.Println("announcing", id, "as", value)
	time.Sleep(30 * time.Second)
	go AnnounceSite(id, site, client)
}

func main() {
	site := NewSite("docker")
	client := etcd.NewClient([]string{os.Getenv("ETCD_URL")})
	go AnnounceSite("docker", site, client)
	select {}
}
