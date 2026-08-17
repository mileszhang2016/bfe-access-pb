# BFE Access Log Protobuf 协议定义

> 文件：`bfe_access_pb/bfe_access.proto`  
> 语法：`proto2`  
> 包名：`bfe_access_pb`

## 1. 顶层消息：BfeLog

`BfeLog` 是所有访问日志的统一容器。

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `product` | `ProductID` | 1 | 产品标识，固定为 `BFE = 18` |
| `timestamp` | `uint64` | 2 | 请求或会话结束时的 Unix 时间戳（秒） |
| `logid` | `uint64` | 3 | 请求日志 ID / 会话日志 ID |
| `log_tag` | `string` | 26 | 日志标签：`req_<product>`、`req_err_<product>`、`session_bfe` 等 |
| `log_type` | `BfeLogType` | 200 | 日志类型：`Request = 0` / `Session = 1` |
| `request_log` | `RequestLog` | 210 | 请求日志详情（`log_type = Request` 时填充） |
| `session_log` | `SessionLog` | 211 | 会话日志详情（`log_type = Session` 时填充） |

### 1.1 日志类型枚举

```protobuf
enum BfeLogType {
    Request = 0;  // HTTP 请求日志
    Session = 1;  // HTTP 会话/连接日志
}
```

## 2. RequestLog（请求日志）

描述一次完整的 HTTP 请求生命周期，字段按语义分区编号。

### 2.1 基础信息（编号 1 - 40）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `err_code` | `string` | 3 | 错误码，成功请求为空字符串 |
| `err_msg` | `string` | 4 | 附加错误信息 |
| `req_header_len` | `uint32` | 5 | 请求头长度 |
| `req_body_len` | `uint32` | 6 | 请求体长度（暂未实现） |
| `session_id` | `uint64` | 7 | 请求所属会话 ID |
| `addr_info` | `ConnAddrInfo` | 8 | 连接地址信息 |
| `client_ip` | `uint32` | 9 | 真实客户端源 IP |
| `req_num` | `uint32` | 11 | 会话内请求序号，从 1 开始 |
| `client_network` | `NetType` | 13 | 客户端 IP 网络类型：`IPv4 = 0` / `IPv6 = 1` |
| `client_ip6` | `string` | 14 | 当客户端为 IPv6 时填充 |

### 2.2 请求头信息（编号 41 - 60）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `proto` | `string` | 41 | 协议：`http1.0`、`http1.1`、`https`、`spdy`、`http2` |
| `header_host` | `string` | 42 | HOST 头 |
| `origin_uri` | `string` | 43 | 原始 URI |
| `final_uri` | `string` | 44 | 重写/解码后的最终 URI，未重写时为空 |
| `referrer` | `string` | 45 | Referer 头 |
| `x_forward_for` | `string` | 47 | X-Forwarded-For 头 |
| `accept_language` | `string` | 48 | Accept-Language 头 |
| `authorization` | `string` | 49 | Authorization 头 |
| `user_agent` | `string` | 50 | User-Agent 头 |
| `transfer_encoding` | `string` | 51 | 请求 Transfer-Encoding |
| `method` | `string` | 52 | HTTP 方法：`GET`、`POST`、`PUT` 等 |
| `delegation` | `string` | 53 | 委托域名 |
| `content_type` | `string` | 54 | Content-Type，主要用于安全分析 |
| `req_headers` | `repeated HttpHeader` | 55 | 关注的请求头列表 |
| `uid` | `string` | 57 | 请求头中的 UID |

### 2.3 Cookie 信息（编号 61 - 90）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `cookie` | `string` | 61 | 完整的 Cookie 内容 |

### 2.4 路由定位信息（编号 101 - 150）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `product` | `string` | 101 | 产品名 |
| `cluster` | `string` | 102 | 集群名 |
| `sub_cluster` | `string` | 103 | 子集群名 |
| `backend_info` | `InstanceInfo` | 104 | 后端实例信息，重试时为最后一次重试的后端 |
| `backend_retry` | `uint32` | 105 | 后端重试次数 |

### 2.5 响应信息（编号 151 - 200）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `res_status_code` | `uint32` | 151 | HTTP 响应状态码 |
| `res_header_len` | `uint32` | 152 | 响应头长度 |
| `res_body_len` | `uint32` | 153 | 响应体长度 |
| `res_content_type` | `string` | 154 | 响应 Content-Type |
| `res_location` | `string` | 155 | 3xx 响应的 Location |
| `res_transfer_encoding` | `string` | 156 | 响应 Transfer-Encoding |
| `res_headers` | `repeated HttpHeader` | 158 | 关注的响应头列表 |

### 2.6 耗时信息（编号 201 - 250）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `all_time` | `uint32` | 201 | 请求总耗时（ms）：从开始读请求到完成发送响应 |
| `read_client_time` | `uint32` | 202 | 从客户端读取请求耗时（ms） |
| `cluster_serve_time` | `uint32` | 203 | 集群处理耗时（ms，含重试） |
| `backend_serve_time` | `uint32` | 204 | 后端处理耗时（ms，重试时为最后一次） |
| `write_client_time` | `uint32` | 205 | 向客户端写响应耗时（ms） |
| `session_offset_time` | `uint32` | 206 | 从会话开始到请求结束的偏移时间（ms） |
| `connect_backend_time` | `uint32` | 207 | 连接后端耗时（ms，重试时为最后一次） |
| `proxy_delay_time` | `uint32` | 208 | BFE 代理延迟耗时（ms） |

### 2.7 AI 可观测性字段（编号 701 - 800）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `ai_apikey` | `string` | 701 | 从 Authorization 头提取的 API Key |
| `ai_apikeytags` | `repeated ApikeyTag` | 702 | API Key 附加标签列表 |
| `ai_requested_model` | `string` | 703 | 客户端请求的原始模型名 |
| `ai_mapped_model` | `string` | 704 | 网关实际路由的目标模型名 |
| `ai_stream` | `bool` | 705 | 是否为流式响应 |
| `ai_prompt_tokens` | `int64` | 706 | 输入 Token 数 |
| `ai_output_tokens` | `int64` | 707 | 输出 Token 数 |
| `ai_total_tokens` | `int64` | 708 | 总 Token 消耗 |
| `ai_ttft_us` | `int64` | 709 | 首 Token 延迟 TTFT（微秒），仅流式请求 |
| `ai_tpot_us` | `int64` | 710 | 输出 Token 平均耗时 TPOT（微秒） |
| `ai_rate_limit_hits` | `repeated RateLimitHit` | 711 | 触发的限流策略列表 |
| `ai_auth_reject_reason` | `string` | 712 | 鉴权拒绝原因 |
| `ai_auth_reject_quota_plans` | `repeated string` | 713 | 拒绝时命中的 Quota Plan ID 列表 |

## 3. SessionLog（会话日志）

描述一条 TCP/HTTP 连接的生命周期。

### 3.1 基础信息（编号 1 - 40）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `err_code` | `string` | 2 | 会话错误码 |
| `err_msg` | `string` | 3 | 附加错误信息 |
| `req_num` | `uint32` | 4 | 该连接上服务的请求数 |
| `read_len` | `uint32` | 5 | 从客户端读取的总字节数（仅含头部） |
| `write_len` | `uint32` | 6 | 写入客户端的总字节数（暂未实现） |
| `addrInfo` | `ConnAddrInfo` | 7 | 连接地址信息 |
| `proto` | `string` | 8 | 协议：`http`、`https`、`spdy`、`http2` |
| `product` | `string` | 9 | 产品名 |

### 3.2 时间信息（编号 41 - 50）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `start_time` | `uint64` | 41 | 会话开始的 Unix 时间戳（秒） |
| `all_time` | `uint32` | 42 | 会话总耗时（ms） |
| `rtt` | `uint32` | 43 | 源 IP 到 BFE 的平滑 RTT（ms） |
| `synRtt` | `uint32` | 44 | 四层网关上的 SYN RTT（ms） |

### 3.3 TLS/SSL 信息（编号 101 - 150）

| 字段 | 类型 | 编号 | 说明 |
|------|------|------|------|
| `tls_version` | `uint32` | 101 | TLS/SSL 版本 |
| `cipher_suite` | `uint32` | 102 | 密码套件 |
| `session_resume` | `bool` | 103 | 是否使用 TLS Session Resume |
| `ocsp_staple` | `bool` | 104 | 是否使用 OCSP Stapling |
| `handshake_time` | `uint32` | 105 | TLS 握手耗时（ms） |
| `tls_cert_common_name` | `string` | 106 | 证书 Common Name |
| `tls_cert_chain_id` | `string` | 107 | 证书链 ID |
| `lost_type` | `LostType` | 108 | TLS 握手失败原因分类 |

## 4. 公共消息与枚举

### 4.1 ApikeyTag

```protobuf
message ApikeyTag {
    optional string tagname  = 1;  // 标签名，如 "dep"
    optional string tagvalue = 2;  // 标签值，如 "op"
}
```

### 4.2 RateLimitHit

```protobuf
message RateLimitHit {
    optional string rate_limit_policy_id = 1;  // 限流策略 ID，如 "rlp-0001"
    optional string rate_limit_type      = 2;  // 限流类型：tpm / rpm / concurrency
    repeated string rule_names           = 3;  // 命中的规则名列表，如 ["win2min","win10min"]
}
```

### 4.3 ConnAddrInfo

```protobuf
message ConnAddrInfo {
    required uint32 bfe_ip          = 1;  // 处理请求的 BFE 服务器 IP
    required uint32 sock_src_ip     = 2;  // socket 源 IP
    required bool   is_trust_src_ip = 3;  // socket 源 IP 是否在信任 IP 列表
    optional uint32 vip             = 4;  // 请求目的 VIP
    optional string vip6            = 5;  // IPv6 VIP
}
```

### 4.4 HttpHeader

```protobuf
message HttpHeader {
    required string key   = 1;
    required string value = 2;
}
```

### 4.5 InstanceInfo

```protobuf
message InstanceInfo {
    required uint32 ip_addr = 1;  // IP 地址
    required uint32 port    = 2;  // 端口
}
```

### 4.6 其他枚举

| 枚举 | 值 | 说明 |
|------|----|------|
| `ProductID.BFE` | 18 | 产品标识 |
| `NetType.IPv4` | 0 | IPv4 |
| `NetType.IPv6` | 1 | IPv6 |
| `LostType.Unknown` | 0 | 未知 |
| `LostType.Bfe` | 1 | BFE 端原因 |
| `LostType.NotBfe` | 2 | 非 BFE 端原因 |

## 5. 预留字段说明

为后续扩展，`RequestLog` 中预留了以下字段区间：

| 编号区间 | 用途 |
|----------|------|
| 251 - 400 | 通用预留 |
| 401 - 500 | WAF 预留 |
| 501 - 600 | 模块信息预留 |
| 601 - 610 | 通用预留 |
| 701 - 800 | AI 可观测性字段 |
