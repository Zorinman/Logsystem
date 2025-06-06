[go-callvis官方](https://github.com/ondrajz/go-callvis)
[graphviz官方](https://graphviz.org/download/)

## 使用工具:
1.**go-callvis**：运行指针分析来构建Go程序的调用图，并且使用数据生成点格式的输出
2.**graphviz**:将输出的内容渲染并显示

## 安装(windows)
### 1.go-callvis
#### 在Vscode中打开终端
输入：`go install github.com/ofabry/go-callvis@latest`
#### 配置环境变量
找到go-callvis.exe安装位置并添加到环境变量:
- 查看GOPATH目录:`go env GOPATH`
- 将`GOPATH目录/bin`添加到Windows环境变量-系统变量-Path中


### 2.graphviz
  - 官网安装：https://graphviz.org/download/
  - 我这里选择：`graphviz-12.2.1 (64-bit) EXE installer [sha256]`
  - **添加环境变量**:
   同理将`安装目录\bin`路径添加到Windows环境变量-系统变量-Path中


 ### 重启VScode
 - 终端输入:`go-callvis  -version `确认go-callvis安装成功
   - 终端输入:`dov -v `确认graphviz安装成功

## 快速使用
在项目目录打开终端直接输入：
`go-callvis  main.go ` 可以直接在网页`http://localhost:7878`生成渲染图
在浏览器界面上可以点击每一个调用模块进一步查看模块内的调用
