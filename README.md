## :star: 概述

本项目是由 Golang 编写的IP控制池，同一个 IP 可以交由不同 `domain` 处理

## :fire: 使用

1. 保证代理服务商获取的 `ip` 格式: `<xxx.xxx.xxx.xxx:ip_v4>:<xxx:port>`

2. 在本地或则服务器上运行本程序

3. 使用程序请求

    ```python
    # 192.168.2.211:8000/get?domain=baidu.com
    requests.get("url", params={"domain": "baidu.com"})
    ```

## :robot: 设置

在 `project.toml` 中设置

- [x] 每个域名的IP使用间隔
- [x] 每个IP的存活时间
- [x] IP池大小
- [ ] 不同域名单独设置 IP 的延迟时间

## :warning: 注意

请下载最新版本
