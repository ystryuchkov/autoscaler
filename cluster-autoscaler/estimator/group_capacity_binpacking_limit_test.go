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
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	testprovider "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/test"
)

func newTestNodeGroup(maxSize int, minSize int, targetSize int) cloudprovider.NodeGroup {
	return testprovider.NewTestNodeGroup("test-ng-1", maxSize, minSize, targetSize,
		true, false, "n1-standard-2", nil, nil)
}

func TestNodeGroupsBinpackingLimit(t *testing.T) {
	testCases := []struct {
		name       string
		nodeGroups []cloudprovider.NodeGroup
		want       int
	}{
		{
			name:       "Handles empty node groups",
			nodeGroups: make([]cloudprovider.NodeGroup, 0),
			want:       0,
		},
		{
			name:       "Handles nil",
			nodeGroups: nil,
			want:       0,
		},
		{
			name:       "Negative capacity (maxSize < targetSize) means no binpacking limit",
			nodeGroups: []cloudprovider.NodeGroup{newTestNodeGroup(0, 100, 3)},
			want:       0,
		},
		{
			name: "Computes capacity correctly",
			nodeGroups: []cloudprovider.NodeGroup{
				newTestNodeGroup(10, 0, 5),   // 5 nodes capacity
				newTestNodeGroup(100, 0, 50), // 50 nodes capacity
				newTestNodeGroup(2, 0, 0),    // 2 nodes capacity
				newTestNodeGroup(10, 0, 20),  // no free capacity
			},
			want: 57,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() { NewGroupCapacityBinpackingLimit(tt.nodeGroups).GetNodeLimit() })
			assert.Equal(t, tt.want, NewGroupCapacityBinpackingLimit(tt.nodeGroups).GetNodeLimit(), tt.nodeGroups)
		})
	}
}
