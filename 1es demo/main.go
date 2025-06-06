package main

import (
	"bytes"
	"fmt"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

//这里有两种方式可以插入数据到Elasticsearch中，分别是使用cli.Create和cli.Index方法。

func main() {
	// Create a new Elasticsearch client
	cli, err := elastic.NewClient(elastic.Config{
		Addresses: []string{"http://192.168.219.133:9200"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Elasticsearch client created successfully")
	//插入1.直接使用JSON字符串作为文档内容，使用cli.Create方法将文档添加到索引中

	//bytes.Buffer 用于创建一个新的字节缓冲区，作为文档的内容
	// 这里的文档内容是一个 JSON 字符串，表示一个用户对象
	//bytes.Buffer类型 是一个实现了 io.Reader 接口的字节缓冲区，可以直接传递给 Elasticsearch 的 API,无需再使用strings.NewReader转换
	doc := bytes.NewBufferString(`{
	    "name":jack,
	    "age": 13
	}`)
	// cli.Create.WithPretty()是可选的：使用 WithPretty() 将 Elasticsearch 的响应格式化为易读的 JSON 格式，并返回到 response 中
	//文档ID必须指定,如果指定文档ID在索引中已存在，则会返回 409 错误表示ID冲突
	response, err := cli.Create("user", "2", doc, cli.Create.WithPretty())
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}

//插入2.将一个 person 对象转换为 JSON 字符串，将 JSON 字符串转换为 io.Reader流式数据使用cli.Index方法添加到 Elasticsearch 的索引user中

// 	type person struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}

// 	p1 := person{Name: "John Doe", Age: 30}
// 	data, err := json.Marshal(p1) // 将person对象序列化编码为 JSON 格式的 []byte 字节切片
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = cli.Index(
// 		"user",
// 		strings.NewReader(string(data)), //使用 string(data) 将 JSON 的字节切片转换为字符串，然后通过 strings.NewReader 将字符串转换为 io.Reader，以便按流式接口读取 JSON 数据
// 		// bytes.NewReader(data), // 也可以直接使用 bytes.NewReader(data) 直接将JSON 字节切片转换为 io.Reader
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Document indexed successfully")
// }

//⭐cli.Create和cli.Index的区别：

// 1. cli.Create
// 作用：
// 用于向指定索引中创建一个文档。
// 如果文档的 ID 已存在，则会返回错误（通常是 HTTP 409 错误，表示 ID 冲突）。
// 特点：
// 必须指定文档的 ID。
// 如果文档 ID 已存在，无法覆盖或更新文档。
// 适用场景：
// 当您需要确保文档的唯一性，并且不希望覆盖已有文档时使用。

// 2. cli.Index
// 作用：
// 用于向指定索引中添加或更新文档。
// 如果文档的 ID 已存在，则会覆盖已有文档。
// 特点：
// 文档的 ID 是可选的。如果未指定，Elasticsearch 会自动生成一个唯一的 ID。
// 支持文档的更新（通过覆盖实现）。
// 适用场景：
// 当您需要更新已有文档或不关心文档的 ID 时使用。
