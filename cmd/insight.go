// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

package cmd

import (
	"github.com/kubearmor-client/insight"
	"github.com/spf13/cobra"
)

var insightOptions insight.Options

// insightCmd represents the insight command
var insightCmd = &cobra.Command{
	Use:   "insight",
	Short: "Observe policy from the discovery engine",
	Long:  `Observe policy from the discovery engine`,
	RunE: func(cmd *cobra.Command, args []string) error {
		insight.StopChan = make(chan struct{})
		if err := insight.StartObserver(insightOptions); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(insightCmd)

	insightCmd.Flags().StringVar(&insightOptions.GRPC, "gRPC", "", "gRPC server information")
	insightCmd.Flags().StringVar(&insightOptions.GRPC, "labels", "", "Labels for resources")
	insightCmd.Flags().StringVar(&insightOptions.GRPC, "containername", "", "Filter according to the Container name")
	insightCmd.Flags().StringVar(&insightOptions.GRPC, "clustername", "", "Filter according to the Cluster name")
	insightCmd.Flags().StringVarP(&insightOptions.Namespace, "namespace", "n", "explorer", "Namespace for resources")
	insightCmd.Flags().BoolVar(&insightOptions.JSON, "json", true, "Flag to print alerts and logs in the JSON format")
}
