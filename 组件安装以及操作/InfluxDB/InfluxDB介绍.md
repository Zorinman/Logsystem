**详细内容**：https://blog.csdn.net/qq_44766883/article/details/131511821

### 核心概念
| InfluxDB 2.x+ 名词 | 传统数据库概念  | 说明                                                                  |
| ------------------ | --------------- | --------------------------------------------------------------------- |
| **Organization**   | 租户/工作空间   | 多租户隔离单位，包含多个 Bucket（类似 MySQL 的实例包含多个数据库）    |
| **Bucket**         | 数据库+保留策略 | 合并了 1.x 的 `database` 和 `retention policy`，直接关联数据保留时间  |
| **Measurement**    | 表（Table）     | 时间序列集合（无需预定义结构，写入时自动创建）                        |
| **Point**          | 行（Row）       | 一条时间序列数据（含时间戳+字段值）                                   |
| **Timestamp**      | 主键            | 时间戳（唯一标识，默认按时间排序）                                    |
| **Tags**           | 带索引的列      | 元数据标签（自动索引，用于高效过滤，如 `host=server1`）               |
| **Fields**         | 不带索引的列    | 实际存储的指标值（如 `cpu_usage=85.2`，支持多种数据类型）             |
| **Series**         | 唯一数据序列    | 由 `Measurement` + `Tags` 组合确定的唯一序列（如 `cpu,host=server1`） |
| **Shard**          | 分区/分片       | 按时间范围分片存储的数据块（如按天分区）                              |
| **Flux**           | SQL 查询语言    | 2.x 默认的脚本式查询语言（功能更强，但语法与 SQL 差异大）             |


### 存储引擎：
2.x 版本将存储引擎分解为：

#### Meta Store：存储元数据(采用嵌入式KV存储 BoltDB)

#### TSI (Time Series Index)**：改进的倒排索引系统
倒排索引(Inverted Index)是一种索引数据结构，它建立了从内容到文档的映射关系（与传统索引相反）。在数据库领域，它特别适合快速查找包含特定值的记录。

**正向索引**：
文档 → 包含的词
```
文档1 → {词A, 词B, 词C}
文档2 → {词A, 词D}

```
**倒排索引**：
词 → 包含该词的文档
```
词A → {文档1, 文档2}
词B → {文档1}
词C → {文档1}
词D → {文档2}

```
**例子**：
- **假设有以下4个series**：
```
series1: m1,tag1=hello,tag2=world cpu=0.5
series2: m1,tag1=hello,tag2=china cpu=0.3  
series3: m1,tag1=hi,tag2=world cpu=0.7
series4: m1,tag1=hi,tag2=china cpu=0.2
```

- **倒排索引的组织方式**：
通过measurement
```
measurement:m1 → [series1, series2, series3, series4]
tag1:hello → [series1, series2]  
tag1:hi → [series3, series4]
tag2:world → [series1, series3]
tag2:china → [series2, series4]


查询的执行过程：
SELECT * FROM m1 WHERE tag1='hello'：

直接查找 tag1=hello 的倒排列表 → 立即得到[series1, series2]

无需扫描其他series
```

- **正向索引组织方式**：
全表扫描：必须遍历所有series的元数据进行一一匹配 
```
series1 → {m1, {tag1:hello, tag2:world}, {cpu:0.5}}
series2 → {m1, {tag1:hello, tag2:china}, {cpu:0.3}}
series3 → {m1, {tag1:hi, tag2:world}, {cpu:0.7}} 
series4 → {m1, {tag1:hi, tag2:china}, {cpu:0.2}}

查询过程分析：
当执行查询 SELECT * FROM m1 WHERE tag1='hello'：

全表扫描：必须遍历所有series的元数据

检查series1：tag1=hello → 匹配

检查series2：tag1=hello → 匹配

检查series3：tag1=hi → 不匹配

检查series4：tag1=hi → 不匹配

返回结果：series1和series2
```


#### TSM Storage Engine：实际数据存储


