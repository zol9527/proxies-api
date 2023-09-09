## 内存结构

```
{
    "www.xxx-01.com": {
        "lock"": 0,
        "domain": "www.xxx-01.com",
        "lastuse": "2023-07-13 16:46:06",
        "pool": [
            {
                "ip": "0.0.0.0",
                "port": 123,
                "acquire_time": "2023-07-13 16:46:06",
                "lastuse_time": "2023-07-13 16:46:06"
            },
            {
                "ip": "0.0.0.0",
                "port": 123,
                "acquire_time": "2023-07-13 16:46:06",
                "lastuse_time": "2023-07-13 16:46:06"
            },
            ...
        ],
    },

    "www.xxx-02.com": {
        ...
    },
    ...
}
```
