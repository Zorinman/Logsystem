### 8 Context上下文
https://www.liwenzhou.com/posts/Go/context/


#### 前置知识：父节点和子节点（副本）：
用`context.WithCancel`举例
context.WithCancel 是 Go 中用于创建可取消的 context 的函数，其签名如下：

`func WithCancel(parent Context) (ctx Context, cancel CancelFunc)`

**输入**：一个父 context（如 context.Background()）。

**输出：**

一个子节点 context（ctx），它是父 context 的副本，但带有独立的取消机制。

一个 cancel 函数，调用它会取消这个新 context。

----

**-父节点（Parent Context）**：调用 WithCancel 时传入的 context（例如 context.Background()）。

**-子节点（副本）（ctx Context）**：WithCancel 返回的新 context，它的特点是：

- 继承父 context 的所有值（通过 ctx.Value(key) 能获取父 context 存储的数据）。

- 拥有自己的 Done 通道（和父 context 的 Done 通道独立）。

- 监听两种取消信号：

     父 context 的取消（父的 Done 被关闭）。

    显式调用返回的 cancel 函数。


**父 context 取消时，所有派生的子 context 也会被取消（通过关闭 Done 通道），第一次传入的ctx是根节点不可以取消-如context.Background()**

**示例：**
**以下代码创建的context节点关系如下：**
```c
context.Background()  // 根 context（根节点，不可取消）
       |
       v
     parent            // 父 context（可取消，通过 cancel() 关闭）
       |
       v
     child  //子context（，可取消，通过 cancel() 关闭）

```
当手动parentCancel()关闭父节点通道时，子节点通道也会随之关闭
    

```go
func main() {
    parent, parentCancel := context.WithCancel(context.Background())
    defer parentCancel() // 确保父 context 最终被取消

    child, childCancel := context.WithCancel(parent)
    defer childCancel() // 确保子 context 最终被取消

    // 监听父 context 的取消
    go func() {
        <-parent.Done()
        fmt.Println("父 context 已取消")
    }()

    // 监听子 context 的取消
    go func() {
        <-child.Done()
        fmt.Println("子 context 已取消")
    }()

    // 手动取消父 context
    parentCancel()

    time.Sleep(time.Second)
}
```

#### 补充：关于为什么context.WithValueGo 官方推荐使用自定义类型作为 context key而不是内建类型如 string

##### 错误示范：两个包使用相同的 string key

假设有两个包：`packageA` 和 `packageB`，它们都往 context 中写入一个 key 为 `"userID"` 的值：

**packageA/a.go：**
```go
package packageA

import "context"

func AddUserIDToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, "userID", "A-User")
}

func GetUserID(ctx context.Context) string {
    return ctx.Value("userID").(string)
}
```

**packageB/b.go：**
```go
package packageB

import "context"

func AddUserIDToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, "userID", "B-User") // 相同 key，不同用途
}

func GetUserID(ctx context.Context) string {
    return ctx.Value("userID").(string)
}
```

**主程序 main.go：**
```go
package main

import (
    "context"
    "fmt"

    "example.com/project/packageA"
    "example.com/project/packageB"
)

func main() {
    ctx := context.Background()

    ctx = packageA.AddUserIDToContext(ctx)
    ctx = packageB.AddUserIDToContext(ctx) // 覆盖了 A 的值

    fmt.Println("From A:", packageA.GetUserID(ctx)) // ❌ 实际输出: "B-User"，不是 "A-User"
    fmt.Println("From B:", packageB.GetUserID(ctx)) // ✅ 输出: "B-User"
}
```

##### 问题：
- `packageA` 想存的是它自己的 `"userID"`，但 `packageB` 也用了同样的 key `"userID"`，导致值被覆盖，`packageA` 无法正确获取自己的值。
- 这种情况在大型项目或引入外部库时非常容易发生。

---

##### 正确做法：使用自定义类型避免冲突

**修改 packageA/a.go：**
```go
package packageA

import "context"

type userIDKey struct{} // 独立 key 类型

func AddUserIDToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, userIDKey{}, "A-User")
}

func GetUserID(ctx context.Context) string {
    val := ctx.Value(userIDKey{})
    if val == nil {
        return ""
    }
    return val.(string)
}
```

**修改 packageB/b.go：**
```go
package packageB

import "context"

type userIDKey struct{} // 即使名字相同，类型不一样，也不会冲突

func AddUserIDToContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, userIDKey{}, "B-User")
}

func GetUserID(ctx context.Context) string {
    val := ctx.Value(userIDKey{})
    if val == nil {
        return ""
    }
    return val.(string)
}
```

**主程序 main.go 不变：**
```go
func main() {
    ctx := context.Background()

    ctx = packageA.AddUserIDToContext(ctx)
    ctx = packageB.AddUserIDToContext(ctx)

    fmt.Println("From A:", packageA.GetUserID(ctx)) // ✅ 输出: "A-User"
    fmt.Println("From B:", packageB.GetUserID(ctx)) // ✅ 输出: "B-User"
}
```

##### 总结：
- 使用自定义类型作为 context key 可以避免不同包或模块之间的 key 冲突。
- 即使两个包的 key 名称相同，只要它们的类型不同，就不会互相干扰。
- 这是 Go 官方推荐的最佳实践，尤其在大型项目中非常重要。
