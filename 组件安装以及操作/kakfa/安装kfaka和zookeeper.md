**使用版本**:wurstmeister/kafka:2.12-2.3.0

 ###  **安装zookeeper** 
`docker run -d --name zookeeper -p 2181:2181 -v /etc/localtime:/etc/localtime -v /root/zookeeper/data:/data -v /root/zookeeper/datalog:/datalog zookeeper:3.5.5`
使用3.5.5版本
/conf 存储着ZooKeeper 的配置文件
**持久化**:
这里将容器内
/datalog挂载到了本地的/root/zookeeper/datalog   （zookeeper的数据快照）
/data挂载到了本地的/root/zookeeper/data         （zookeeper的事务日志）



 ### **安装kafka**
 `docker run  -d --name kafka -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=192.168.219.132:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://192.168.219.132:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092  -v /root/kafka/kafkalogs:/kafka -d wurstmeister/kafka:2.12-2.3.0`

 使用2.12-2.3.0版本
` /opt/kafka/config`保存了kafka的配置文件
**持久化**:
这里将容器内
`/kafka`目录挂载到了本地的/root/kafka/kafkalogs  (Kafka 消息日志存储目录，可以在/config配置文件的server.properties 中查找log.dir关键词查看)


### 关系

        ┌──────────┐
        │ Sarama   │← Go语言客户端，用于连接 Kafka
        └──────────┘
              │
              ▼
        ┌──────────┐
        │  Kafka   │← 消息系统（生产者/消费者）
        └──────────┘
              │
              ▼
        ┌────────────┐
        │ ZooKeeper  │← Kafka 的早期依赖，用于kafka集群注册与发现、leader选举和follower信息同步、负载均衡、存储集群元数据（v2.x）
        └────────────┘

**目前由于KIP-500提案，Kafka将逐步去除对ZooKeeper的依赖，转而使用社区自研的基于Raft算法的共识机制来替代zookeeper的功能**