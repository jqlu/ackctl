package client

import "github.com/denverdino/aliyungo/cs"

type Cluster struct {
	*cs.ClusterType

	KubeConfig *cs.ClusterConfig
}
