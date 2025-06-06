[Go语言操作](https://golang.halfiisland.com/community/database/Elasticsearch.html)

[核心概念介绍](https://blog.csdn.net/weixin_42081445/article/details/144748629)

[概念以及REstful接口介绍](https://www.cnblogs.com/hxjcore/p/18182067)
[李文周-Restful接口介绍](https://www.liwenzhou.com/posts/Go/elasticsearch/)
## Elasticsearch vs 关系型数据库(MySQL)核心概念对照表
| Elasticsearch 概念 | 关系型数据库 概念    | 说明                                                                                  |
| ------------------ | -------------------- | ------------------------------------------------------------------------------------- |
| **Index**          | Database             | ES 中的一个 Index 相当于一个数据库，包含多个类型（7.x 起已弃用 Type 概念）            |
| **Document**       | Row（记录）          | Document 是 JSON 格式的记录，类似于关系数据库中的一行数据                             |
| **Field**          | Column（列）         | Document 中的字段，对应表中的列                                                       |
| **Type**（已废弃） | Table（表）          | 早期版本 ES 中 Type 类似表的概念，7.x 后已弃用                                        |
| **Mapping**        | Schema（表结构）     | 定义文档中字段的类型及其索引方式                                                      |
| **Shard**          | 数据分片             | 分片就是将一个逻辑数据库拆成多个物理片段，类似 MySQL 的分库分表，避免单一硬件资源不足 |
| **Replica**        | 副本（备份）         | Shard 的冗余副本，提高可用性和容错性                                                  |
| **Query DSL**      | SQL（查询语言）      | Elasticsearch 的查询语言，基于 JSON 构建                                              |
| **Inverted Index** | B+ Tree / Hash Index | 倒排索引是全文搜索的核心，区别于传统数据库的索引结构                                  |


