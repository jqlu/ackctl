package client

import (
	"encoding/json"
	"fmt"
	"github.com/denverdino/aliyungo/common"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"time"
)

type Taints []v1.Taint
type SpotPrice struct {
	InstanceType string `json:"instance_type"`
	PriceLimit   string `json:"price_limit"`
}

type TagItemType struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Tags []*TagItemType

type NodePoolInfo struct {
	NodePoolId      string        `json:"nodepool_id"`
	RegionId        common.Region `json:"region_id"`
	Name            string        `json:"name"`
	Created         time.Time     `json:"created"`
	Updated         time.Time     `json:"updated"`
	IsDefault       bool          `json:"is_default"`
	NodepoolType    string        `json:"type"`
	ResourceGroupId string        `json:"resource_group_id"`
}

type NodePoolStatus struct {
	TotalNodes    int    `json:"total_nodes"`
	OfflineNodes  int    `json:"offline_nodes"`
	ServingNodes  int    `json:"serving_nodes"`
	RemovingNodes int    `json:"removing_nodes"`
	FailedNodes   int    `json:"failed_nodes"`
	InitialNodes  int    `json:"initial_nodes"`
	HealthyNodes  int    `json:"healthy_nodes"`
	State         string `json:"state"`
}

type BasicNodePool struct {
	NodePoolInfo   `json:"nodepool_info"`
	NodePoolStatus `json:"status"`
}

type NodeKubernetesConfig struct {
	CpuPolicy      string `json:"cpu_policy"`
	Runtime        string `json:"runtime,omitempty"`
	RuntimeVersion string `json:"runtime_version"`

	Labels Labels `json:"labels"`
	Taints Taints `json:"taints"`
}

type Labels []Label

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ScalingGroup struct {
	ScalingGroupId string `json:"scaling_group_id"`
	ScalingGroupConfig
}

type ScalingGroupConfig struct {
	VSwitches          []string   `json:"vswitch_ids"`
	SecurityGroupId    string     `json:"security_group_id"`
	InstanceTypes      []string   `json:"instance_types"`
	SystemDiskCategory string     `json:"system_disk_category"`
	SystemDiskSize     int64      `json:"system_disk_size"`
	DataDisks          []DataDisk `json:"data_disks"`

	Platform string `json:"platform"`
	ImageId  string `json:"image_id"`

	LoginPassword string `json:"login_password"`
	KeyPair       string `json:"key_pair"`

	InstanceChargeType string      `json:"instance_charge_type"`
	Period             int         `json:"period"`
	PeriodUnit         string      `json:"period_unit"`
	AutoRenew          bool        `json:"auto_renew"`
	AutoRenewPeriod    int         `json:"auto_renew_period"`
	SpotStrategy       string      `json:"spot_strategy"`
	SpotPriceLimit     []SpotPrice `json:"spot_price_limit"`

	Tags Tags `json:"tags"`
}

type DataDisk struct {
	Size     int    `json:"size"`
	Category string `json:"category"`
}

type AutoScaling struct {
	Enable           bool   `json:"enable"`
	ScalingGroupType string `json:"type"`
	MaxInstances     int64  `json:"max_instances"`
	MinInstances     int64  `json:"min_instances"`
}

type TEEConfig struct {
	TEEType   string `json:"tee_type"`
	TEEEnable bool   `json:"tee_enable"`
}

type NodePool struct {
	BasicNodePool
	NodeKubernetesConfig `json:"kubernetes_config"`
	ScalingGroup         `json:"scaling_group"`
	AutoScaling          `json:"auto_scaling"`
	TEEConfig            `json:"tee_config"`
}

type NodePoolList struct {
	NodePools []*NodePool `json:"nodepools"`
}

type NodePoolScaleRequest struct {
	Count int `json:"count"`
}

type Node struct {
	InstanceId     string `json:"instance_id"`
	InstanceName   string `json:"instance_name"`
	NodeName       string `json:"node_name"`
	NodeStatus     string `json:"node_status"`
	State          string `json:"state"`
	InstanceRole   string `json:"instance_role"`
	InstanceType   string `json:"instance_type"`
	InstanceStatus string `json:"instance_status"`
}

type PageInfo struct {
	TotalCount int `json:"total_count"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type NodeListResponse struct {
	Nodes []Node   `json:"nodes"`
	Page  PageInfo `json:"page"`
}

type RemoveNodeRequest struct {
	Nodes       []string `json:"nodes"`
	ReleaseNode bool     `json:"release_node"`
	DrainNode   bool     `json:"drain_node"'`
}

func ListNodePools(clusterId string) ([]*NodePool, error) {
	client := GetSdkClient()

	var nodePoolList NodePoolList
	path := fmt.Sprintf("/clusters/%s/nodepools", clusterId)
	err := client.Request("GET", path, nil, nil, &nodePoolList)
	if err != nil {
		return nil, err
	}

	return nodePoolList.NodePools, nil
}

func ListNodes(clusterId, nodePoolId string) ([]*Node, error) {
	client := GetSdkClient()
	nodes := make([]*Node, 0)

	var response NodeListResponse
	path := fmt.Sprintf("/clusters/%s/nodes", clusterId)
	query := map[string]string{
		"nodepool_id": nodePoolId,
		"pageNumber":  "1",
		"pageSize":    "100",
		"state":       "all",
	}

	for {
		err := client.Request("GET", path, query, nil, &response)
		if err != nil {
			return nil, err
		}

		for _, n := range response.Nodes {
			nodes = append(nodes, &n)
		}

		if len(nodes) >= response.Page.TotalCount {
			break
		}

		query["pageNumber"] = strconv.Itoa(response.Page.PageNumber + 1)
	}

	return nodes, nil
}

func ScaleNodePool(clusterId, nodePoolId string, count int) error {
	client := GetSdkClient()

	path := fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, nodePoolId)
	request := NodePoolScaleRequest{Count: count}
	body, _ := json.Marshal(request)

	err := client.Request("PUT", path, nil, &body, nil)
	if err != nil {
		return err
	}

	return nil
}

func RemoveNode(clusterId, nodeName string, release, drain bool) error {
	client := GetSdkClient()

	path := fmt.Sprintf("/api/v2/clusters/%s/nodes/remove", clusterId)
	request := RemoveNodeRequest{
		Nodes: []string{nodeName},
		ReleaseNode: release,
		DrainNode: drain,
	}
	body, _ := json.Marshal(request)

	err := client.Request("POST", path, nil, &body, nil)
	if err != nil {
		return err
	}

	return nil
}
