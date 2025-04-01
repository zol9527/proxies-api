# 代理IP池管理系统

## :star: 概述

本项目是由 Golang 编写的IP控制池，同一个 IP 可以交由不同 `domain` 处理。系统提供高性能的代理IP管理，支持自动获取、清理和分配IP，并针对不同域名提供独立的IP轮询机制和延迟设置。

## :rocket: 特性

- 🔄 IP自动轮询分配与负载均衡
- 🕒 可配置的IP生命周期管理
- 🌐 支持多域名独立管理
- 📊 基于域名的请求频率限制
- 📥 同时支持从URL和本地文件获取IP
- 🧹 自动清理过期IP
- 🛡️ 内置健康检查机制

## :hammer_and_wrench: 安装

```bash
# 克隆仓库
git clone https://github.com/zol9527/proxies-api.git
cd proxies-api

# 安装依赖
go mod download

# 编译项目
go build -o proxy-pool

# 运行服务
./proxy-pool
```

## :fire: 使用

1. 保证代理服务商获取的 `ip` 格式: `<xxx.xxx.xxx.xxx:ip_v4>:<xxx:port>`

2. 在本地或服务器上运行本程序

3. 使用程序请求

    ```python
    # 示例：请求 http://localhost:8000/proxy?domain=baidu.com
    import requests
    
    response = requests.get("http://localhost:8000/proxy", params={
        "domain": "baidu.com",
        "delay": "2,5"  # 可选参数，设置延迟范围
    })
    
    if response.status_code == 200:
        data = response.json()
        proxy = f"{data['ip']}:{data['port']}"
        print(f"获取到代理: {proxy}")
    else:
        print(f"错误: {response.json()['msg']}")
    ```

## :gear: API 文档

### 获取代理 IP

**请求**:
- 路径: `/proxy`
- 方法: `GET`
- 参数:
  - `domain` (必填): 请求域名
  - `delay` (可选): 延迟设置，可以是逗号分隔的数组 (如 `2,5`) 或多个参数 (如 `delay=2&delay=5`)

**成功响应** (200):
```json
{
  "ip": "192.168.1.1",
  "port": 8080,
  "delay": [2, 5]
}
```

**错误响应** (400/429/500):
```json
{
  "code": 400,
  "msg": "❌ 缺少必要的domain参数"
}
```

### 限流机制

系统内置请求限流机制，防止单个域名请求过多，默认限制为每秒10个请求，突发上限为20个请求。

## :robot: 配置

在 `config/proxy_pool_config.toml` 中设置：

```toml
[server]
port = 8000  # HTTP服务端口

# 代理池配置
[pool]
size = 100  # IP池大小

# 代理源配置
[proxy_source]
urls = [  # 代理源URL列表
    "http://example.com/proxy1",
    "http://example.com/proxy2"
]
request_freq = 10  # URL请求频率(秒)
local_path = [  # 本地代理文件路径
    "/var/proxy_pool/proxy.txt",
    "/var/proxy_pool/proxy2.txt"
]
read_freq = 10  # 文件读取频率(秒)

# IP配置
[ip]
lifetime_seconds = 900  # IP生存时间(秒)
delay = [0, 5]  # 默认延迟范围(秒)
```

配置项说明:
- ✅ 每个域名的IP使用间隔
- ✅ 每个IP的存活时间
- ✅ IP池大小
- ✅ 不同域名单独设置 IP 的延迟时间 (通过API参数)

## :wrench: 开发

项目结构:
```
proxies-api/
├── cmd/            # 命令行入口
├── config/         # 配置文件
├── internal/       # 内部实现
│   ├── api/        # HTTP接口
│   ├── config/     # 配置解析
│   ├── log/        # 日志工具
│   ├── models/     # 数据结构
│   └── pool/       # IP池管理
├── main.go         # 主程序入口
└── README.md       # 文档
```

### 单元测试

```bash
go test ./...
```

## :warning: 注意事项

- 确保配置文件中的路径有效
- 对于高并发场景，建议适当调整限流参数
- 系统会定期清理过期IP，无需手动操作

## :page_facing_up: 许可证

本项目采用 MIT 许可证
