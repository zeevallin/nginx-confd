package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/samalba/dockerclient"
)

func newTlsConfig(path string) *tls.Config {

	ca, _ := ioutil.ReadFile(fmt.Sprintf("%s/ca.pem", path))

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)

	key := fmt.Sprintf("%s/key.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)

	pair, err := tls.LoadX509KeyPair(cert, key)
	logFatalIf(err)

	config := tls.Config{
		InsecureSkipVerify: (DOCKER_TLS_VERIFY != "1"),
		ClientAuth:         tls.RequireAnyClientCert,
		Certificates:       []tls.Certificate{pair},
		RootCAs:            pool,
	}

	return &config
}

func newDocker() dockerclient.Client {
	var (
		err    error
		docker dockerclient.Client
	)

	if DOCKER_CERT_PATH == "" {
		docker, err = dockerclient.NewDockerClient(DOCKER_HOST, nil)
	} else {
		config := newTlsConfig(DOCKER_CERT_PATH)
		docker, err = dockerclient.NewDockerClient(DOCKER_HOST, config)
	}

	logFatalIf(err)

	return docker
}

func newMappings() []Mapping {
	var mappings []Mapping
	for _, s := range strings.Split(MAPPING, ",") {
		mappings = append(mappings, newMapping(s))
	}
	return mappings
}

func newMapping(arg string) Mapping {
	s := strings.Split(arg, ":")
	p, _ := strconv.Atoi(s[2])
	m := Mapping{
		Upstream: s[0],
		Pattern:  s[1],
		Port:     p,
	}
	return m
}

func newServer(host string, port int) Server {
	return Server{
		Host: fmt.Sprintf("%s:%v", host, port),
		Config: map[string]string{
			"fail_timeout": "0",
		},
	}
}

func isLinkedName(name string) bool {
	linked, _ := regexp.MatchString("/", name[1:])
	return linked
}

func isMappedName(mapping Mapping, name string) bool {
	matched, _ := filepath.Match(mapping.Pattern, name[1:])
	return matched
}

func isMappedPort(m Mapping, port dockerclient.Port) bool {
	return (m.Port == port.PrivatePort)
}

func handleSignals() {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)

	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, os.Kill)

	select {
	case <-sigInt:
		log.Println("Received SIGINT")
		os.Exit(0)
	case <-sigKill:
		log.Println("Received SIGKILL")
		os.Exit(0)
	}
}

func announce(key string, server Server, etc *etcd.Client) {
	marshal, _ := json.Marshal(server)
	value := string(marshal)
	log.Println("Announcing", server.Host, "as", value)
	etc.Set(key, value, 10)
}

func checkContainers(mappings []Mapping, docker dockerclient.Client, etc *etcd.Client) {

	for {

		containers, err := docker.ListContainers(false, true, "")
		logFatalIf(err)

		pointers := map[string]Server{}

		for _, container := range containers {

			var name string

			for _, containerName := range container.Names {
				if !isLinkedName(containerName) {
					name = containerName
					break
				}
			}

			for _, mapping := range mappings {

				var ports []int

				if isMappedName(mapping, name) {
					for _, port := range container.Ports {
						if isMappedPort(mapping, port) {
							ports = append(ports, port.PublicPort)
						}
					}
				}

				for _, port := range ports {
					key := fmt.Sprintf("/nginx/servers/%s/%s/%s", CLUSTER, mapping.Upstream, container.Id[0:12])
					pointers[key] = newServer(HOST, port)
				}
			}
		}

		for key, server := range pointers {
			go announce(key, server, etc)
		}

		time.Sleep(5 * time.Second)
	}

}

func run() {

	mappings := newMappings()
	docker := newDocker()
	etc := etcd.NewClient([]string{ETCD_URL})

	go checkContainers(mappings, docker, etc)

	log.Println("Running...")
}

func main() {

	captureEnvironment()

	go handleSignals()
	go run()

	select {}
}
