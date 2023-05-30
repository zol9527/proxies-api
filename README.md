## :star: 概述

本项目是由 Golang 编写的IP控制池，同一个 IP 可以交由不同 `domain` 处理

## :fire: 使用

```python
# 192.168.2.211:8000/get?domain=baidu.com
requests.get("url", params={"domain": "baidu.com"})
```

## :robot: 设置

在 `project.toml` 中设置

- [x] 每个域名的IP使用间隔
- [x] 每个IP的存活时间
- [x] IP池大小
