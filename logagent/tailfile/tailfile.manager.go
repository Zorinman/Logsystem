package tailfile

import (
	"fmt"
	"logagent/common"
)

//tailTask的管理者

type tailTaskManager struct {
	tailTaskMap      map[string]*tailTask       // 存储 tailTask 实例的映射
	collectEntryList []common.CollectEntry      // 存储日志收集配置项的列表
	confChan         chan []common.CollectEntry // 用于接收配置项的通道

}

var ttMgr *tailTaskManager // 声明一个全局的 tailTaskManager 实例
func Init(allConf []common.CollectEntry) (err error) {
	// 创建一个新的 tailTaskManager 实例
	ttMgr = &tailTaskManager{
		tailTaskMap:      make(map[string]*tailTask, 50),   // 初始化 tailTaskMap
		collectEntryList: allConf,                          // 初始化 collectEntryList
		confChan:         make(chan []common.CollectEntry), // 初始化 confChan
	}
	//allConf存了若干日志收集项
	//针对每一个日志收集项，创建一个 tailTask 实例
	//每个新建的tailTask示例存储到 ttMgr.tailTaskMap 中，方便后续管理
	// 并启动一个 goroutine 来运行 tailTask 实例
	for _, conf := range allConf {
		tt, err := newTailTask(conf.Path, conf.Topic) // 创建一个新的 tailTask 实例
		if err != nil {
			fmt.Printf("路径%s文件tailTask追踪实例创建失败:%s", conf.Path, err)
			return err // 如果创建失败，返回错误
		}
		fmt.Println("为conf.Path创建tailTask追踪实例成功")    // 打印成功消息
		ttMgr.tailTaskMap[tt.path] = tt              // 将 新创建的tailTask 实例存储到 tailTaskMap 中 方便后续管理
		go tt.Run()                                  // 启动 tailTask 实例，开始追踪日志文件
		fmt.Printf("路径%s文件Run监听线程创建成功\n", conf.Path) // 打印成功消息

	}
	go ttMgr.WatchConf() // 启动一个 goroutine 来监听日志配置项的变化

	return
}

// etcd中的WatchConf监听etcd中日志项配置的变化 → 如果有变化调用SendNewConf()函数发送到confChan通道 →  ttMgr.WatchConf() 函数从confChan接收新配置项开始处理
func (ttMgr *tailTaskManager) WatchConf() {
	for {
		newConf := <-ttMgr.confChan                      // //等待etcd的WatchConf通知有新的日志收集项
		fmt.Printf("Tail收到新的日志收集配置项%v，开始处理...", newConf) // 打印收到新配置的消息
		for _, conf := range newConf {
			// 1.如果添加了新日志配置项：为新的日志收集配置项创建新的 tailTask 实例，已经存在的日志收集配置项保持不变
			if ttMgr.isExist(conf) {
				continue // 如果 ttMgr.tailTaskMap 中已经存在该日志收集配置项，则跳过
			}
			// 不存在则创建新的 tailTask 实例
			tt, err := newTailTask(conf.Path, conf.Topic) // 创建一个新的 tailTask 实例
			if err != nil {
				fmt.Printf("新的路径%s文件tailTask追踪实例创建失败:%s", conf.Path, err)
				continue
			}
			fmt.Printf("新的路径%s文件tailTask追踪实例创建成功:", conf.Path) // 打印成功消息
			ttMgr.tailTaskMap[tt.path] = tt                    // 将 新创建的tailTask 实例存储到 tailTaskMap 中 方便后续管理
			go tt.Run()                                        // 启动 tailTask 实例，开始追踪日志文件
			fmt.Printf("新的路径%s文件Run监听线程创建成功\n", conf.Path)
			// 打印成功消息
		}
		// 2.如果删除了旧的日志配置项：从 ttMgr.tailTaskMap 中删除对应的 tailTask 实例，并停止tail.TailFile创建的文件追踪实例以及关闭对应的go协程
		for key, task := range ttMgr.tailTaskMap {
			exist := false // 标记是否存在
			for _, conf := range newConf {
				if key == conf.Path { // 如果当前 key 在新的配置项中存在
					exist = true // 标记为存在
					break        // 跳出循环
				}
			}
			if !exist { // 如果当前 key 在新的配置项中不存在
				//  Tail的Stop方法可以用来关闭tail.TailFile创建的文件追踪实例tt.Obj
				// 由于tt.Obj关闭，tt.Obj.Lines通道也会关闭，从而go tt.Run()创建的处理该tailTask 实例的协程也会跳出for循环关闭
				task.tObj.Stop()
				delete(ttMgr.tailTaskMap, key)                              // 从 tailTaskMap 中删除对应的 tailTask 实例
				fmt.Printf("路径%s文件tailTask追踪实例已删除且对应goroutine协程已停止\n", key) // 打印删除成功消息
			}
		}

	}
}

// 判断ttMgr.tailTaskMap中是否有对应的日志收集配置项
func (ttMgr *tailTaskManager) isExist(conf common.CollectEntry) bool {
	// 检查 ttMgr.tailTaskMap 中是否存在指定的日志收集配置项
	_, exists := ttMgr.tailTaskMap[conf.Path] // 检查 tailTaskMap 中是否存在指定路径的 tailTask 实例
	return exists                             // 返回是否存在的结果
}

func SendNewConf(newConf []common.CollectEntry) {
	//接收新的日志收集配置项
	//将新的日志收集配置项发送到 confChan 通道
	ttMgr.confChan <- newConf // 将新的配置项发送到 confChan 通道
	fmt.Println("有新配置项发送到 confChan 通道")
}
