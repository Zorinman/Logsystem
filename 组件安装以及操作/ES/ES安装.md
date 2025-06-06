**参考官方安装脚本文件中的yaml格式**
https://www.elastic.co/docs/solutions/search/run-elasticsearch-locally

**脚本文件下载**：https://elastic.co/start-local

⭐请直接使用yaml文件,从以下内容中复制粘贴到vim打开的文件中会格式会打乱
### docker compose文件：
```yml
version: "3.9"

services:
  elasticsearch:
    image: elasticsearch:7.10.1
    container_name: elasticsearch
    volumes:
      - /root/es/data:/usr/share/elasticsearch/data
    ports:
      - 9200:9200
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
      - xpack.security.enabled=false #（默认关闭）
    #因为xpack插件需要付费，trial只是免费30天使用所以直接关闭
    #   - xpack.license.self_generated.type=trial
    #   - ELASTIC_PASSWORD=123
    #   - xpack.security.enabled=true #开启Elasticsearch 的安全模块（需要用户名密码登录)
    #   - xpack.security.http.ssl.enabled=false #关闭 HTTP 层的 SSL 加密(默认也是关闭)
    #   - xpack.ml.use_auto_machine_memory_percent=true
    #注意单节点情况下healthcheck不要设置为要求green，因为单节点集群状态是 YELLOW，永远不会变成 green，导致 Docker 认为容器“不健康”启动失败。
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9200"]
      interval: 10s
      timeout: 5s
      retries: 10


  kibana:
    depends_on:
      elasticsearch:
        condition: service_healthy #健康检查，Es顺利启动之后才会启动kibana
    image: kibana:7.10.1
    container_name: kibana
    volumes:
      - /root/es/kibana/data:/usr/share/kibana/data
    #   - ./config/telemetry.yml:/usr/share/kibana/config/telemetry.yml
    ports:
      - 5601:5601
    environment:
      - SERVER_NAME=kibana
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200 #设置 Kibana 要连接的 Elasticsearch 地址
    #   - ELASTICSEARCH_USERNAME=elastic
    #   - ELASTICSEARCH_PASSWORD=123
      - ELASTICSEARCH_PUBLICBASEURL=http://localhost:9200 #设置当 Kibana 向外暴露时显示的 Elasticsearch 地址
    healthcheck:
      test:
        [
          "CMD-SHELL",
          "curl -s -I http://kibana:5601 | grep -q 'HTTP/1.1 302 Found'",
        ]
      interval: 10s
      timeout: 10s
      retries: 30

```