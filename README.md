# ackctl

`ackctl` 是阿里云容器服务 Kubernetes 版（ACK）的命令行管理工具。

## 安装

macOS: https://ackctl.oss-cn-hangzhou.aliyuncs.com/macOS/ackctl

Linux: https://ackctl.oss-cn-hangzhou.aliyuncs.com/linux/ackctl

## 配置
### 使用阿里云 CLI 的配置
`ackctl` 支持直接读取 AK 模式的阿里云 CLI 配置，无需重新配置，可以直接使用。其他模式暂不支持。

### 手动配置
```bash
# ackctl configure  
? Access Key Id: ****************
? Access Key Secret: ******************************
Ackctl configured.
```

## 使用

### 获取集群列表
```bash
ackctl get cluster
ID                                State         Region         Type               Name
ccd9ca1cebbad43e7b61815dd25a4712e running       ap-southeast-1 ManagedKubernetes  ttttt
c8ae985fb582e4b438f785626420891c2 running       cn-hangzhou    ManagedKubernetes  hangzhou-production
c10a26830ab474634bd5098664eae4021 running       cn-hangzhou    ManagedKubernetes  hangzhou-testing
c3af8a161483d4f8a9253a1200b8dcfcc running       cn-shenzhen    ManagedKubernetes  app-center-testing
c40e56f7f896749479b8e0feec7c37014 running       cn-hangzhou    ManagedKubernetes  ol-edge
c982785cc8bb64a47a066016ab2f4dcf5 waiting       cn-zhangjiakou ExternalKubernetes ol-ex
ce7bc8633abe84520bc365328acdcfe6d running       cn-shanghai    Ask                ol-ask-v2
cf53ee0a037dd41ccbddc2efc68779cee running       cn-beijing     Ask                ol-ask-v1
```

### 创建集群
```bash
# ackctl create cluster -f testdata/create.json
Starting to create cluster c548b7a5c0a8c4b0a8b7ef71b173a8789
```

### 配置集群 kubeconfig
可以单独将一个集群的 kubeconfig 配置到本地：
```bash
# ackctl use cluster c3af8a16148
? Config /Users/jonas/.kube/config exists, overwrite? Yes
/Users/jonas/.kube/config updated to use cluster
```

使用`--all`参数，将多个集群的 kubeconfig 合并后配置到本地，然后使用`kubectl config use-context <context-name>`即可在集群之间切换，context 的名称为集群名称。

```bash
# ackctl use cluster --all
? Config /Users/jonas/.kube/config exists, overwrite? Yes
Merged kubeConfigs of 7 clusters into: /Users/jonas/.kube/config.
Use 'kubectl config get-contexts' to list contexts.
Use 'kubectl config use-context' to select context.

# kubectl config get-contexts
CURRENT   NAME                  CLUSTER                             AUTHINFO                                             NAMESPACE
          app-center-testing    c3af8a161483d4f8a9253a1200b8dcfcc   c3af8a161483d4f8a9253a1200b8dcfcc-kubernetes-admin
          hangzhou-production   c8ae985fb582e4b438f785626420891c2   c8ae985fb582e4b438f785626420891c2-kubernetes-admin
          hangzhou-testing      c10a26830ab474634bd5098664eae4021   c10a26830ab474634bd5098664eae4021-kubernetes-admin
          ol-ask-v1             cf53ee0a037dd41ccbddc2efc68779cee   cf53ee0a037dd41ccbddc2efc68779cee-kubernetes-admin
          ol-ask-v2             ce7bc8633abe84520bc365328acdcfe6d   ce7bc8633abe84520bc365328acdcfe6d-kubernetes-admin
          ol-edge               c40e56f7f896749479b8e0feec7c37014   c40e56f7f896749479b8e0feec7c37014-kubernetes-admin

# kubectl config use-context app-center-testing
Switched to context "app-center-testing".
```

### 删除集群
```bash
# ackctl delete cluster c9c86da441c
? Are you sure to delete cluster test(c9c86da441ced45208c85a9ff2eca5b4b)? Cluster cannot be restored after deletion Yes
Starting to delete cluster test(c9c86da441ced45208c85a9ff2eca5b4b)
```

### 查询集群节点池
`ackctl get nodepool`，使用`--cluster-id`或`-c`指定集群 ID。
```bash
# ackctl get nodepool -c c8ae985fb582e4b438f785626420891c2
Name             Id                                 State  Total Serving Offline
default-nodepool np7081d5cc325e4019a6f740ee1ca7c7da active 3     3       0
nodepool1        npdf426eef3bc54957bf0a59e314a15fa6 active 1     1       0
```

### 查询节点列表
`ackctl get nodepool`，使用`--cluster-id`或`-c`指定集群 ID；使用`--node-pool-id`或`-p`指定节点池。

```bash
# ackctl get node -c c8ae985fb582e4b438f785626420891c2 -p np7081d5cc325e4019a6f740ee1ca7c7da
Instance Id            Node Name               Instance Status Role   Instance Type   Node Status
i-bp1g91typs6213lznpj1 cn-hangzhou.10.1.41.111 Ready           Worker ecs.hfc6.xlarge running
i-bp1g91typs6213lznpj1 cn-hangzhou.10.1.41.111 Ready           Worker ecs.hfc6.xlarge running
i-bp1g91typs6213lznpj1 cn-hangzhou.10.1.41.111 Ready           Worker ecs.hfc6.xlarge running
```

### 扩容节点池
`ackctl scale nodepool <node-pool-id>`，使用`--cluster-id`或`-c`指定集群 ID；使用`--increment`指定扩容数量。
```bash
# ackctl scale nodepool np7081d5cc325e4019a6f740ee1ca7c7da -c c8ae985fb582e4b438f785626420891c2 --increment 1
Staring to scale node pool np7081d5cc325e4019a6f740ee1ca7c7da of cluster c8ae985fb582e4b438f785626420891c2

# ackctl get nodepool -c c8ae985fb582e4b438f785626420891c2
Name             Id                                 State   Total Serving Offline
default-nodepool np7081d5cc325e4019a6f740ee1ca7c7da scaling 3     3       0
nodepool1        npdf426eef3bc54957bf0a59e314a15fa6 active  1     1       0
```

## Feature roadmap
- 移除节点
- 创建/删除节点池