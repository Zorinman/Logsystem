package main

// tail 包来实现对日志文件的实时监控和读取
// 这个示例代码展示了如何使用 tail 包来监控一个日志文件，并实时读取新添加的日志行
import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	filename := "test.log" // 指定要监控的日志文件
	config := tail.Config{
		ReOpen:    true,                                 // 如果文件被轮转（如日志切割），重新打开文件
		Follow:    true,                                 // 跟随文件的增长（类似于 tail -f）
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件末尾开始读取
		Poll:      true,                                 // 使用轮询模式检测文件变化
		MustExist: true,                                 // 确保文件存在
	}
	tails, err := tail.TailFile(filename, config) // 使用指定的配置开始追踪文件
	if err != nil {
		fmt.Println("文件追踪出错:", err)
		return
	}

	// var (
	// 	msg *tail.Line // 声明一个变量用于存储从文件读取的行
	// 	ok  bool       // 声明一个布尔值用于检查行是否成功读取
	// )
	for {
		msg, ok := <-tails.Lines // 从 tail 中读取一行
		if !ok {
			fmt.Printf("文件关闭后重新打开, 文件名:%s\n", tails.Filename)
			time.Sleep(time.Second * 1) // 等待 1 秒后重试
			continue                    // 继续下一次循环读取下一行
		}
		fmt.Printf("读取到的行: %s\n", msg.Text) // 打印从文件读取的行
	}
}
