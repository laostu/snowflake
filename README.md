snowflake
====
[![GoDoc](https://godoc.org/github.com/bwmarrin/snowflake?status.svg)](https://godoc.org/github.com/bwmarrin/snowflake) [![Go report](http://goreportcard.com/badge/bwmarrin/snowflake)](http://goreportcard.com/report/bwmarrin/snowflake) [![Coverage](http://gocover.io/_badge/github.com/bwmarrin/snowflake)](https://gocover.io/github.com/bwmarrin/snowflake) [![Build Status](https://travis-ci.org/bwmarrin/snowflake.svg?branch=master)](https://travis-ci.org/bwmarrin/snowflake) [![Discord Gophers](https://img.shields.io/badge/Discord%20Gophers-%23info-blue.svg)](https://discord.gg/0f1SbxBZjYq9jLBk)

snowflake 是一个 [Go](https://golang.org/) 语言包，提供以下功能：
* 一个非常简洁的雪花算法 ID 生成器。
* 解析已有雪花 ID 的方法。
* 将雪花 ID 转换为多种其他数据类型以及反向转换的方法。
* JSON Marshal/Unmarshal 函数，方便在 JSON API 中使用雪花 ID。
* 单调时钟计算，防止时钟回拨。

**如需有关此包或 Go 语言一般性讨论的帮助，请加入 [Discord Gophers](https://discord.gg/0f1SbxBZjYq9jLBk) 聊天服务器。**

## 状态
此包已趋于稳定并基本完成。未来的新增功能将尽量避免修改已有函数的 API。

### ID 格式
默认情况下，ID 格式遵循原始的 Twitter 雪花算法格式。
* ID 整体是一个存储在 int64 中的 63 位整数
* 41 位用于存储毫秒精度的时间戳，使用自定义纪元（epoch）
* 10 位用于存储节点 ID - 范围为 0 到 1023
* 12 位用于存储序列号 - 范围为 0 到 4095

### 自定义格式
你可以通过设置 `snowflake.NodeBits` 和 `snowflake.StepBits` 值来改变节点 ID 和序列号所使用的位数。请注意，这两个值共享最多 22 位可用位。你不必使用全部 22 位。

### 自定义纪元
默认情况下，此包使用 纪元：1767225600000（即 2026年01月01日 00:00:00）。你可以通过将 `snowflake.Epoch` 设置为以毫秒为单位的时间值来自定义纪元。

### 自定义注意事项
设置自定义纪元或位数值时，需要在使用 snowflake 包中的任何函数（包括 `NewNode()`）之前进行设置。否则，自定义的值将无法正确生效。

### 工作原理
每次生成 ID 时，工作流程如下：
* 先添加 NodeID。
* 然后添加序列号，从 0 开始，同一毫秒内每生成一个 ID 递增一次。如果在同一毫秒内生成的 ID 数量导致序列号溢出，生成函数将暂停直到下一毫秒。
* 最后后在后续位中添41位存储毫秒精度的时间戳。

默认的 ID 格式如下所示：
```
+--------------------------------------------------------------------------+
| 1 Bit 未使用  |  10 Bit 节点ID  |   12 Bit 序列号  | 41 Bit 时间戳 |
+--------------------------------------------------------------------------+
```

使用默认设置时，每个节点 ID 每毫秒最多可生成 4096 个唯一 ID。

## 快速开始

### 安装

前提是你已经有一个可用的 Go 开发环境，如果没有，请先参考[此页面](https://golang.org/doc/install)。

```sh
go get github.com/laostu/snowflake
```


### 使用

将包导入到你的项目中，然后使用唯一的节点编号构造一个新的 snowflake Node。默认设置允许节点编号范围为 0 到 1023。如果你设置了自定义的 NodeBits 值，则需要自行计算节点编号的范围。通过节点对象调用 `Generate()` 方法来生成并返回唯一的雪花 ID。

请注意，你创建的每个节点必须具有唯一的节点编号，即使跨多个服务器也是如此。如果节点编号不唯一，生成器将无法保证在所有节点之间生成唯一的 ID。

**示例程序：**

```go
package main

import (
	"fmt"

	"github.com/laostu/snowflake"
)

func main() {

	// 创建一个节点编号为 1 的新 Node
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 生成一个雪花 ID
	id := node.Generate()

	// 以几种不同方式打印该 ID
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// 打印该 ID 的时间戳
	fmt.Printf("ID Time  : %d\n", id.Time())

	// 打印该 ID 的节点编号
	fmt.Printf("ID Node  : %d\n", id.Node())

	// 打印该 ID 的序列号
	fmt.Printf("ID Step  : %d\n", id.Step())

  // 一步生成并打印
  fmt.Printf("ID       : %d\n", node.Generate().Int64())
}
```

### 性能

使用默认设置时，该雪花算法生成器在大多数系统上足够快，每毫秒可生成 4096 个唯一 ID，这也是雪花 ID 格式支持的最大值。即每次操作大约耗时 243-244 纳秒。

由于雪花生成器是单线程的，主要限制将取决于你系统上单个处理器的最高速度。

要对生成器进行基准测试，请在 snowflake 包目录下运行以下命令：

```sh
go test -run=^$ -bench=.
```