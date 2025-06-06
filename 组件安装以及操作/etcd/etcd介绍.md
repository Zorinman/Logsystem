## etcd使用了两种端口:
### 2379 端口（Client 通信端口）
**用途：**

用于 客户端与 etcd 集群的通信（如 etcdctl、应用程序 API 调用）。

处理所有客户端的读写请求（Put/Get/Delete 等操作）。

**通信协议：**

HTTP/HTTPS（默认明文，生产环境建议用 TLS 加密）。

gRPC（etcd v3 API 默认使用，高性能二进制协议）。

**典型场景：**

```
#客户端连接示例
`etcdctl --endpoints=http://192.168.1.100:2379 get /key`
`curl http://192.168.1.100:2379/v3/kv/range -X POST -d '{"key": "Zm9v"}'  # v3 API`
```
###  2380 端口（Peer 通信端口）
**用途：**

用于 etcd 集群节点间的内部通信（如选举、数据同步、心跳检测）。

仅在集群模式下使用（单节点模式可不开放此端口）。

**通信协议：**

HTTP/HTTPS（同 2379，但仅用于节点间交互）。

Raft 协议（基于 HTTP 的日志复制和领导者选举）。

```
# 集群节点配置示例（每个节点的 --initial-advertise-peer-urls 需包含 2380 端口）
etcd --name node1 \
     --listen-peer-urls http://0.0.0.0:2380 \
     --listen-client-urls http://0.0.0.0:2379 \
     --initial-advertise-peer-urls http://192.168.1.100:2380 \
     --initial-cluster "node1=http://192.168.1.100:2380,node2=http://192.168.1.101:2380"

```
### 问题:二进制协议更加高效,为什么内部通信使用gRPC
### 答：
**早期设计选择：**
etcd 最初基于 Raft 协议 的实现（如 CoreOS 的 raft 库）**直接使用了 HTTP/1.x**，因为 Raft 的日志复制和心跳机制本身是简单的请求-响应模型，HTTP 已足够。

向后兼容：
保持 HTTP 可以避免因协议升级导致的集群分裂风险（例如混合版本集群中部分节点不支持 gRPC）。

2. **Raft 协议的特性决定**
简单性优先：
Raft 的日志复制和领导者选举本质是顺序化的消息交换（**非高并发场景**），HTTP/1.x 的延迟对性能影响有限。

流量模式：
**节点间通信主要是低频控制消息（心跳、日志同步），而非高频数据流**。二进制协议的优势在大规模数据传输中更明显，而 etcd 的 Raft 消息通常很小。

## 常用使用命令

**存储键值对：**

`etcdctl --endpoints=http://localhost:2379 put mykey "myvalue"`
**读取键值对：**

`etcdctl --endpoints=http://localhost:2379 get mykey`
**列出所有键值对：**

`etcdctl --endpoints=http://localhost:2379 get --prefix ""`
**删除键值对：**

`etcdctl --endpoints=http://localhost:2379 del mykey`
**查看 etcd 集群状态：**

`etcdctl --endpoints=http://localhost:2379 endpoint status`
**备份 etcd 数据：**

`etcdctl --endpoints=http://localhost:2379 snapshot save backup.db`
**恢复 etcd 数据：**

`etcdctl snapshot restore backup.db --data-dir /path/to/etcd/data-dir`