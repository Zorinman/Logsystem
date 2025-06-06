package main

import (
	"fmt"
	"net"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") //向外部地址（Google 的 DNS 服务器）创建一个 UDP“伪连接”（不会实际发送数据）。
	if err != nil {
		fmt.Println("获取外网 IP 地址失败:", err)
		return "" // 如果创建连接失败，返回空字符串
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr) // 获取本地端的地址（IP + 端口) 并通过类型断言将接口net.Addr转换为 *net.UDPAddr 类型
	// 使用类型断言的原因：
	//conn.LocalAddr()返回的是一个接口类型 net.Addr，这是一个抽象接口，代表“某种网络地址”。你不能直接访问 IP、Port 之类的字段，因为接口只定义了行为（方法），没有字段
	//而实现这个接口的具体结构体 *net.UDPAddr是其中之一，*net.UDPAddr结构体实现了net.Addr接口定义的 Network() 和 String() 方法
	// 通过类型断言将接口类型转换为具体的结构体类型，这样就可以访问Ip、Port等字段了
	fmt.Println(localAddr.String())
	return localAddr.IP.String() // 返回本机的外网 IP 地址
	// 注意：这种方法可能不适用于所有网络环境，特别是在 NAT 或代理服务器后面时。
}

func main() {
	ip := GetOutboundIP() // 调用函数获取外网 IP 地址
	if ip != "" {
		fmt.Println("本机的外网 IP 地址是:", ip) // 打印获取到的外网 IP 地址
	} else {
		fmt.Println("未能获取到外网 IP 地址") // 如果获取失败，打印提示信息
	}
}
