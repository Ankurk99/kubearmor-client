// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

//go:build insight
// +build insight

package insight

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// import protobuf from policy discovery

	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	GRPC string
	JSON bool
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

	// create a client
	conn, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := NewObservabilityClient(conn)

	// var response opb.Response
	response := client.SysObservabilityData(context.Background(), &Data.Request)

	log.Printf("%v", response)

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	close(StopChan)

	return nil
}
