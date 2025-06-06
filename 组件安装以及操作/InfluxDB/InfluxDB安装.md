
**版本**：influxdb:2.4
**使用docker compose启动**:
**参考官方模板**:https://docs.influxdb.org.cn/influxdb/v2/install/use-docker-compose/
```yaml
version: '3'

services:
  docker-influxdb:
    image: influxdb:2.4
    container_name: influxdb
    restart: always
    ports:
      - "8086:8086" #HTTP UI and API port
    environment:
      DOCKER_INFLUXDB_INIT_MODE: "setup"
      DOCKER_INFLUXDB_INIT_USERNAME: "root" #创建管理员用户
      DOCKER_INFLUXDB_INIT_PASSWORD: "a123456a" #创建管理员密码，太简单会报错
      DOCKER_INFLUXDB_INIT_ORG: "logagent" #组织名称
      DOCKER_INFLUXDB_INIT_BUCKET: "my-bucket"
    volumes:
      - "/root/influxDB/data:/var/lib/influxdb2"
      - "/root/influxDB/conf:/etc/influxdb2"
```

**这里 `DOCKER_INFLUXDB_INIT_ADMIN_TOKEN` 可以不用指定**：
首次启动时，InfluxDB 会自动完成初始化：

创建管理员用户 root。

生成一个 默认的超级用户令牌（admin token），拥有所有权限。

创建组织 chudaozhe 和存储桶 my-bucket。

你不需要手动设置 influxdb2-admin-token，因为初始化流程会自动生成它

**查询 admin token**:
- 1.`docker logs influxdb`
- 2.通过 InfluxDB UI  进入 Load Your Data > API Tokens 查看

