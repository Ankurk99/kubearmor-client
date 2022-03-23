//go:build cilium
// +build cilium

package vm

import (
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

// NetworkPolicyRequest is the request type used for sending the Cilium
// network policy to KVM service
type NetworkPolicyRequest struct {
	Type   string                 `json:"type"`
	Object v2.CiliumNetworkPolicy `json:"object"`
}
