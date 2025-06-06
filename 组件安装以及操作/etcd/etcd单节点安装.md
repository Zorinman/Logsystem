**参考**；https://cloud.tencent.com/developer/article/2392231
**使用版本**:bitnami/etcd:3.5.12
注意这里使用的是bitnami的etcd镜像，因为bitnami的有bash能够进行容器交互
Bitnami 的 etcd 镜像 默认开启身份验证，如果你没有设置 ETCD_ROOT_PASSWORD，容器就会拒绝启动。
所有在environment中要添加以下内容
      `- ALLOW_NONE_AUTHENTICATION=yes`

⭐每次重启虚拟机可能导致etcd写入键值对失败，可以尝试重启docker解决:`systemctl restart docker`
## 方式1：docker run 快速部署(没有持久化)
这里 由于是单节点所有不需要过多设置
```shell
 docker run -d \
  --name etcd-single-node \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  -e ETCD_ADVERTISE_CLIENT_URLS=http://192.168.219.129:2379 \
  -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
  -e ETCD_NAME=etcd0 \
  -v /root/etcd:/bitnami/etcd \
  bitnami/etcd:3.5.12

```

 
## 方式2：docker-compose 部署
### 这里使用docker-compose.yaml进行安装

### 编写docker-compose.yaml文件

    `image`也可以自己去docker hub 上找需要安装的版本

```yaml
version: '3'

services:
  etcd:
    container_name: etcd1
    image: bitnami/etcd:3.5.21
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd1
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://192.168.219.129:2379
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://192.168.219.129:2380
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://192.168.219.129:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
      - ETCD_LOGGER=zap
      - ETCD_LOG_LEVEL=info
    volumes:
      - /root/etcd/data:/bitnami/etcd
      - "/etc/localtime:/etc/localtime:ro"
    ports:
      - 2379:2379
      - 2380:2380
    restart: always


```
#### 在yaml文件位置执行 docker compose up -d 启动


## 常用的环境参数解析
| 参数                               | 含义                                                                 | 说明                                                           |
| ---------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------- |
| `ETCD_LISTEN_CLIENT_URLS`          | etcd 监听客户端请求的地址(通常为:http://0.0.0.0:2379)                | 监听所有网络接口的 2379 端口，客户端通过该端口访问 etcd。      |
| `ETCD_ADVERTISE_CLIENT_URLS`       | etcd 对外广播的客户端访问地址 (通常为http://宿主机ip:2379)           | 告诉客户端和其他服务访问 etcd 时使用的地址，通常应为实际 IP。  |
| `ETCD_LISTEN_PEER_URLS`            | etcd 监听集群内节点间通信的地址(通常为:http://0.0.0.0:2380)          | 监听所有接口的 2380 端口，用于节点之间的内部通信。             |
| `ETCD_INITIAL_ADVERTISE_PEER_URLS` | 节点启动时向其他节点广告自己的 peer 地址(通常为http://宿主机ip:2380) | 其他节点通过此地址访问本节点的 peer 服务，通常为实际 IP 地址。 |
| `ETCD_INITIAL_CLUSTER_TOKEN`       | 集群唯一标识符                                                       | 区分不同 etcd 集群，防止集群混淆，默认 `etcd-cluster`。        |
| `ETCD_INITIAL_CLUSTER`             | 定义初始集群成员及对应 peer 地址(通常为http://宿主机ip:2380)         | 格式 `name=peerURL`，多节点用逗号分隔，指定集群成员列表。      |
| `ETCD_INITIAL_CLUSTER_STATE`       | 集群状态                                                             | `new` 表示新集群启动，`existing` 表示加入已有集群。            |
