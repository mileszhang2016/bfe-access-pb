# BFE Access Log 总体说明与 b2log 设计

## 1. 概述

`bfe-access-pb` 仓库用于保存 BFE（Baidu Front End）访问日志相关的 Protocol Buffers 定义以及二进制日志操作库。核心组成：

- `bfe_access_pb/bfe_access.proto`：访问日志的 Protobuf Schema 定义（详见 [protobuf.md](./protobuf.md)）。
- `bfe_access_pb/bfe_access.pb.go`：由 `protoc` 生成的 Go 代码。
- `b2log`：BFE access to log 的读写库，负责将 Protobuf 序列化后的数据封装成带魔数和头部的二进制记录。

## 2. 整体架构

BFE Access Log 采用两层结构：

```
+---------------------------------------------------+
| B2Log Record（二进制存储层）                       |
|  +---------------------------------------------+  |
|  | Header：魔数 + 版本 + 长度 + 时间戳           |  |
|  +---------------------------------------------+  |
|  | Payload：Protobuf 序列化后的 BfeLog 消息       |  |
|  +---------------------------------------------+  |
+---------------------------------------------------+
```

- **日志语义层**：`BfeLog` 消息统一描述一次 HTTP 请求（`RequestLog`）或一次 HTTP 连接/会话（`SessionLog`）。
- **二进制存储层**：`b2log` 将序列化后的 `BfeLog` 负载封装为自描述的二进制记录，便于落盘、传输与后续解析。

## 3. B2Log 二进制存储格式

> 位置：`b2log` 包

B2Log 记录由 **Header + Payload** 两部分组成。

### 3.1 Header 结构

| 字段 | 类型 | 字节数 | 说明 |
|------|------|--------|------|
| `MagicNumber` | `uint32` | 4 | 魔数 `0xB0AEBEA7`，小端序 |
| `Version` | `uint32` | 4 | 版本号，当前为 `1` |
| `UnCompressLen` | `uint32` | 4 | 未压缩 Payload 长度 |
| `CompressLen` | `uint32` | 4 | 压缩后 Payload 长度，当前未压缩时为 `0` |
| `TimeStamp` | `uint64` | 8 | 日志生成时间戳 |

- **Header 总长度**：24 字节（`HEADER_SIZE`）。
- **字节序**：Little Endian。
- **最大单条记录长度**：100 KB（`MAX_RECORD_LEN`）。

### 3.2 Record 结构

```
+--------+---------+
| Header | Payload |
+--------+---------+
  24字节   N字节
```

- 当 `CompressLen == 0` 时，Payload 长度为 `UnCompressLen`，内容为序列化后的 `BfeLog` protobuf 消息。
- 当 `CompressLen != 0` 时，当前实现直接跳过该记录，并返回 `ErrCompressed` 错误。

## 4. 读写接口说明

### 4.1 读取接口

#### `BuffParse(buffer []byte) ([]Record, []byte)`

从给定的字节缓冲区中解析出所有完整的 B2Log 记录。

- **入参**：`buffer` 待解析的二进制数据。
- **返回值**：
  - `records`：解析出的记录列表，每条记录为 Payload 字节。
  - `buffer`：解析后剩余的不完整数据，等待下一次数据补充后继续解析。
- **容错**：遇到魔数错误时会通过 `tryFindNextStart` 尝试同步到下一个合法记录头。

### 4.2 写入接口

#### `HeaderWrite(buffer []byte, payloadLen int) error`

向 `buffer` 写入 B2Log Header。

- **入参**：
  - `buffer`：长度至少为 `HEADER_SIZE` 的字节数组。
  - `payloadLen`：Payload 长度，写入到 `UnCompressLen` 字段。
- **行为**：自动填充魔数、版本号以及当前时间戳；`CompressLen` 固定为 `0`。

## 5. 使用示例

### 5.1 写入日志

```go
import "github.com/bfenetworks/bfe-access-pb/b2log"

payload := []byte("serialized BfeLog protobuf bytes")
buff := make([]byte, b2log.HEADER_SIZE + len(payload))

// 写入头部
if err := b2log.HeaderWrite(buff, len(payload)); err != nil {
    // handle error
}

// 写入 Payload
copy(buff[b2log.HEADER_SIZE:], payload)

// 写入文件或网络，省略...
```

### 5.2 读取日志

```go
import "github.com/bfenetworks/bfe-access-pb/b2log"

data, err := ioutil.ReadFile("pb_access.log")
if err != nil {
    // handle error
}

records, remaining := b2log.BuffParse(data)
for _, record := range records {
    // record 为 Payload，可进一步反序列化为 BfeLog
}
```

## 6. 代码生成

运行仓库根目录下的 `build.sh` 可从 `.proto` 重新生成 `.pb.go`：

```sh
sh build.sh
```

当前使用 `protoc-gen-go@v1.35.0`，要求 Go 1.22 环境。
