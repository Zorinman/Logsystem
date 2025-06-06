### 问：为什么`producer`创建的`msg.Value`类型是 `sarama.Encoder`，`consumer`消费消息时msg.Value类型是`[]byte`?

### 答：
其实是producer最终将msg发送到kafka后msg.Value才转为 []byte。

#### 这样设计的原因：
之发送时才转为 []byte而不是一开始就是[]byte，是因为`sarama.Encoder`是一个接口类型，你可以传入任何自定义数据类型，只要它实现了 Encode() 方法就行，这样就能更灵活、更通用地支持多种数据类型。而**最终所有类型的数据在计算机内存或存储中的本质表示都是字节（即 []byte，里面存放着二进制数据（0 和 1））**
