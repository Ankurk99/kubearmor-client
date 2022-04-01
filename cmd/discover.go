// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

package cmd

import (
	"github.com/kubearmor/kubearmor-client/discover"
	"github.com/spf13/cobra"
)

var discoverOptions discover.Options

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover applicable policies",
	Long:  `Discover applicable policies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := discover.DiscoverPolicy(discoverOptions); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().StringVar(&discoverOptions.GRPC, "gRPC", "", "gRPC server information")
	discoverCmd.Flags().StringVarP(&discoverOptions.Policy, "policy", "p", "", "Default: cilium and kubeArmor")
	//discoverCmd.Flags().StringVarP(&discoverOptions.filter, "namespace", "n", "", "Namespace for resources")
}
