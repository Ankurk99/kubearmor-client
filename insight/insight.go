// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

//go:build insight
// +build insight

package insight

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	opb "github.com/kubearmor-client/insight/protobuf"
	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	GRPC          string
	Clustername   string
	Labels        string
	Containername string
	Namespace     string
	JSON          bool
}

// StopChan Channel
var StopChan chan struct{}

// GetOSSigChannel Function
func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}

// StartObserver Function
func StartObserver(o Options) error {
	gRPC := ""

	if o.GRPC != "" {
		gRPC = o.GRPC
	} else {
		if val, ok := os.LookupEnv("DISCOVERY_SERVICE"); ok {
			gRPC = val
		} else {
			gRPC = "localhost:9089"
		}
	}

	fmt.Println("gRPC server: " + gRPC)

	data := &opb.Data{
		Request:   "observe",
		Namespace: o.Namespace,
	}

	// create a client
	conn, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	// listen for interrupt signals

	client := opb.NewObservabilityClient(conn)

	// var response opb.Response
	response, err := client.SysObservabilityData(context.Background(), data)

	// log.Printf("%v\n", response)
	str := ""
	arr, _ := json.Marshal(response)
	str = fmt.Sprintf("%s\n", string(arr))

	log.Printf("%s\n", str)

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	close(StopChan)

	return nil
}
