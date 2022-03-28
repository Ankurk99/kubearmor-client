// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

//go:build insight
// +build insight

package insight

import (
	"context"
	"encoding/json"
	"errors"
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
	Labels        string
	Containername string
	Clustername   string
	Fromsource    string
	Namespace     string
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

func StartInsight(o Options) error {
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
		Request:       "observe",
		Labels:        o.Labels,
		ContainerName: o.Containername,
		ClusterName:   o.Clustername,
		FromSource:    o.Fromsource,
		Namespace:     o.Namespace,
	}

	// create a client
	conn, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		return errors.New("Could not connect to the server. Possible troubleshooting:\n- Check if discovery engine is running\n- Create a portforward to discovery engine service using\n\t\033[1mkubectl port-forward -n explorer service/knoxautopolicy --address 0.0.0.0 --address :: 9089:9089\033[0m\n- Configure grpc server information using\n\t\033[1mkarmor log --grpc <info>\033[0m")
	}
	defer conn.Close()

	client := opb.NewObservabilityClient(conn)

	// var response opb.Response
	response, err := client.SysObservabilityData(context.Background(), data)

	str := ""
	arr, _ := json.MarshalIndent(response, "", "    ")
	str = fmt.Sprintf("%s\n", string(arr))

	log.Printf("%s \n", str)

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	close(StopChan)

	return nil
}
