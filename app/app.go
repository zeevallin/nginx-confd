package main
import (
  "os"
  "fmt"
  "time"
  "log"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/coreos/go-etcd/etcd"
  "github.com/samalba/dockerclient"
)

type Server struct {
  Host string `json:"host"`
  ServerConfig map[string]string `json:"config,omitempty"`
}

func PanicIf(err error) {
  if err != nil {
    panic(err)
  }
}

func announce(k string, h string, e *etcd.Client) {
  docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
  container, _ := docker.InspectContainer(h)
  ip := os.Getenv("DOCKER_HOST")

  for _, b := range container.NetworkSettings.Ports["3000/tcp"] {
    host := fmt.Sprintf("%s:%s", ip, b.HostPort)
    server := &Server{
      Host: host,
      ServerConfig: map[string]string{"fail_timeout":"0"},
    }
    m, _ := json.Marshal(server)
    v := string(m)

    log.Println("announcing", host, "as", v)
    e.Set(k, v, 3)
  }

  time.Sleep(1500 * time.Millisecond)
  go announce(k, h, e)
}

func main() {
  h, _ := os.Hostname()
  c := os.Getenv("CLUSTER")
  u := "app"
  e := etcd.NewClient([]string{os.Getenv("ETCD_URL")})

  k := fmt.Sprintf("/nginx/servers/%s/%s/%s", c, u, h)
  go announce(k, h, e)

  m := martini.Classic()
  m.Get("/*", func() string {
    return fmt.Sprintf("%s | %s | %s", h, c, time.Now())
  })
  m.Run()
}