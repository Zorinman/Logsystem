# Kafka 常见操作命令

## 1. 启动 Kafka 和 Zookeeper

### 启动 Zookeeper（docker部署容器启动则代表Zookeeper启动）
```bash
zookeeper-server-start.sh config/zookeeper.properties
```

### 启动 Kafka（docker部署容器启动则代表kafka启动）
```bash
kafka-server-start.sh config/server.properties
```

---

docker 部署后直接使用`docker logs -f` 即可进行实时日志消息监听

## 2. 创建主题
```bash
kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic <topic-name>
```
- `--bootstrap-server`：Kafka 集群地址。
- `--replication-factor`：副本因子。
- `--partitions`：分区数。
- `--topic`：主题名称。

---

## 3. 列出所有主题
```bash
kafka-topics.sh --list --bootstrap-server localhost:9092
```

---

## 4. 查看主题详情
```bash
kafka-topics.sh --describe --bootstrap-server localhost:9092 --topic <topic-name>
```

---

## 5. 删除主题
```bash
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic <topic-name>
```

> 注意：删除主题需要 Kafka 开启 `delete.topic.enable=true` 配置。

---

## 6. 生产消息
```bash
kafka-console-producer.sh --broker-list localhost:9092 --topic <topic-name>
```
- 输入消息后按回车发送。

---

## 7. 消费消息（如果没有提前创建消费组组，则每执行一次消费会自动创建消费者组）
```bash
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic <topic-name> --from-beginning
```
- `--from-beginning`：从头开始消费消息。

---

## 8. 检查消费者组
### 列出消费者组
```bash
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list
```

### 查看消费者组详情
```bash
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group <group-name>
```

---

## 9. 删除消费者组
```bash
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --delete --group <group-name>
```

---

## 10. 测试生产者和消费者
### 启动生产者
```bash
kafka-console-producer.sh --broker-list localhost:9092 --topic <topic-name>
```

### 启动消费者
```bash
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic <topic-name> --from-beginning
```