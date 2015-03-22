package main

import (
	"os"
)

var (
	DOCKER_HOST       string // Unix socket or HTTP address for the docker host machine
	DOCKER_TLS_VERIFY string // Boolean to say if docker client should verify TLS
	DOCKER_CERT_PATH  string // Path to directory containing all docker keys and certificates
	ETCD_URL          string // Host name and protocol of the etcd server
	MAPPING           string // Comma separated list to determine what containers to announce
	CLUSTER           string // Name of the load balanced cluster
	HOST              string // Address to the host docker machine
)

// Capture environment variables and assign them to constants
func captureEnvironment() {
	DOCKER_HOST = os.Getenv("DOCKER_HOST")
	DOCKER_TLS_VERIFY = os.Getenv("DOCKER_TLS_VERIFY")
	DOCKER_CERT_PATH = os.Getenv("DOCKER_CERT_PATH")
	ETCD_URL = os.Getenv("ETCD_URL")
	MAPPING = os.Getenv("MAPPING")
	CLUSTER = os.Getenv("CLUSTER")
	HOST = os.Getenv("HOST")

	if HOST == "" {
		HOST = "127.0.0.1"
	}

	if CLUSTER == "" {
		CLUSTER = "local"
	}
}
