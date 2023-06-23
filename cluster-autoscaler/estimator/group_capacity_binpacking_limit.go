/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package estimator

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/klog/v2"
)

type groupCapacityBinpackingLimit struct {
	nodeGroups []cloudprovider.NodeGroup
}

// GetNodeLimit returns total available capacity in all node groups given to this limit.
// It assumes that node groups are valid (initialized and have target size defined)
func (l *groupCapacityBinpackingLimit) GetNodeLimit() int {
	var totalCapacity int
	for _, nodeGroup := range l.nodeGroups {
		nodeGroupTargetSize, err := nodeGroup.TargetSize()
		// Should not ever happen as only valid node groups are passed to estimator
		if err != nil {
			klog.Errorf("Error while computing available capacity of node group %v: can't get target size of the group", nodeGroup.Id(), err)
			continue
		}
		groupCapacity := nodeGroup.MaxSize() - nodeGroupTargetSize
		if groupCapacity > 0 {
			totalCapacity += groupCapacity
		}
	}
	return totalCapacity
}

// NewGroupCapacityBinpackingLimit returns a node count threshold for EstimationLimiter
// based on total available capacity in all given node groups. It expects that
// invalid node groups (not initialized, with no target size) are already filtered out.
func NewGroupCapacityBinpackingLimit(nodeGroups []cloudprovider.NodeGroup) BinpackingLimit {
	return &groupCapacityBinpackingLimit{
		nodeGroups: nodeGroups,
	}
}
