package client

import (
	"fmt"
	"github.com/denverdino/aliyungo/cs"
	"strings"
)

type ClusterCreateResponse struct {
	ClusterId string `json:"cluster_id"`
}

func FindClusterByPrefix(prefix string) (*cs.ClusterType, error) {
	client := GetCsClient()
	clusters, err := client.DescribeClusters("")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cluster list")
	}

	var filtered []cs.ClusterType
	for _, c := range clusters {
		if strings.HasPrefix(c.ClusterID, prefix) {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no cluster matching the prefix %s", prefix)
	}

	if len(filtered) > 1 {
		return nil, fmt.Errorf("ambiguous cluster id prefix: %s", prefix)
	}

	return &filtered[0], nil
}

func GetAccessibleClusters() ([]*cs.ClusterType, error) {
	client := GetCsClient()
	clusters, err := client.DescribeClusters("")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cluster list")
	}

	var result []*cs.ClusterType
	for i, c := range clusters {
		if c.ClusterType != "aliyun" && (c.State == cs.Running || c.State == cs.Scaling) {
			result = append(result, &clusters[i])
		}
	}

	return result, nil
}

func CreateCluster(params []byte) (*ClusterCreateResponse, error) {
	client := GetSdkClient()

	var response ClusterCreateResponse
	err := client.Request("POST", "/clusters", nil, &params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
