package client

import (
	"fmt"
	"github.com/denverdino/aliyungo/cs"
	"github.com/imdario/mergo"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
	"sync"
)

func GetMergedKubeConfig(clusters []*cs.ClusterType) ([]byte, error) {
	var result []*Cluster
	wg := sync.WaitGroup{}
	ch := make(chan *Cluster, len(clusters))
	client := GetCsClient()

	for _, c := range clusters {
		wg.Add(1)
		go func(c *cs.ClusterType) {
			defer wg.Done()
			config, err := client.DescribeClusterUserConfig(c.ClusterID, false)
			if err != nil {
				log.Printf("failed to get kubeconfig for cluster %s: %v", c.Name, err)
			}
			ch <- &Cluster{
				ClusterType: c,
				KubeConfig:  config,
			}
		}(c)
	}

	wg.Wait()
	close(ch)

	for v := range ch {
		result = append(result, v)
	}

	if len(result) == 1 {
		return []byte(result[0].KubeConfig.Config), nil
	}

	return mergeKubeConfigs(result)
}

func mergeKubeConfigs(clusters []*Cluster) ([]byte, error) {
	configs := make([]*clientcmdapi.Config, 0)
	for _, c := range clusters {
		if config, err := clientcmd.Load([]byte(c.KubeConfig.Config)); err == nil {
			updateConfig(c, config)
			configs = append(configs, config)
		}
	}

	mapConfig := clientcmdapi.NewConfig()

	for _, kubeconfig := range configs {
		mergo.Merge(mapConfig, kubeconfig, mergo.WithOverride)
	}

	// merge all of the struct values in the reverse order so that priority is given correctly
	// errors are not added to the list the second time
	nonMapConfig := clientcmdapi.NewConfig()
	for i := len(configs) - 1; i >= 0; i-- {
		kubeconfig := configs[i]
		mergo.Merge(nonMapConfig, kubeconfig, mergo.WithOverride)
	}

	// since values are overwritten, but maps values are not, we can merge the non-map config on top of the map config and
	// get the values we expect.
	config := clientcmdapi.NewConfig()
	mergo.Merge(config, mapConfig, mergo.WithOverride)
	mergo.Merge(config, nonMapConfig, mergo.WithOverride)

	return clientcmd.Write(*config)
}

func updateConfig(cluster *Cluster, config *clientcmdapi.Config) {
	clusterKeys := make([]string, 0)
	for k, _ := range config.Clusters {
		clusterKeys = append(clusterKeys, k)
	}
	for _, k := range clusterKeys {
		config.Clusters[cluster.ClusterID] = config.Clusters[k]
		delete(config.Clusters, k)
	}

	authInfoKeys := make([]string, 0)
	authInfoName := ""
	for k, _ := range config.AuthInfos {
		authInfoKeys = append(authInfoKeys, k)
	}
	for _, k := range authInfoKeys {
		authInfoName = fmt.Sprintf("%s-%s", cluster.ClusterID, k)
		config.AuthInfos[authInfoName] = config.AuthInfos[k]
		delete(config.AuthInfos, k)
	}

	contextKeys := make([]string, 0)
	for k, v := range config.Contexts {
		contextKeys = append(contextKeys, k)
		v.Cluster = cluster.ClusterID
		v.AuthInfo = authInfoName
	}

	for _, k := range contextKeys {
		config.Contexts[cluster.Name] = config.Contexts[k]
		delete(config.Contexts, k)
	}
}
