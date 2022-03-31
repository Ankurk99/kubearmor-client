// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

package discover

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	opb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/worker"
	types "github.com/accuknox/auto-policy-discovery/src/types"
	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	Policy string
	GRPC   string
}

func ConvertPolicy(o Options) error {
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

	data := &opb.WorkerRequest{
		Policytype: o.Policy,
	}

	conn, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := opb.NewWorkerClient(conn)

	var response *opb.WorkerResponse
	response, err = client.Convert(context.Background(), data)
	if err != nil {
		return errors.New("could not connect to the server. Possible troubleshooting:\n- Check if discovery engine is running\n- Create a portforward to discovery engine service using\n\t\033[1mkubectl port-forward -n explorer service/knoxautopolicy --address 0.0.0.0 --address :: 9089:9089\033[0m\n- Configure grpc server information using\n\t\033[1mkarmor log --grpc <info>\033[0m")
	}

	policy := []types.KnoxNetworkPolicy{}
	err = json.Unmarshal(response.CiliumPolicy, &policy)

	str := ""
	arr, _ := json.MarshalIndent(response, "", "    ")
	str = fmt.Sprintf("%s\n", string(arr))

	log.Printf("policy: %v \n", policy)
	log.Printf("response: %v \n", str)

	return nil
}

func DiscoverPolicy(o Options) error {

	if o.Policy == "cilium" {
		o.Policy = "network"
	} else if o.Policy == "kubearmor" {
		o.Policy = "system"
	} else {
		log.Println("Policy type not recognized")
	}

	ConvertPolicy(o)

	return nil
}
