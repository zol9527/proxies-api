This file is a merged representation of the entire codebase, combined into a single document by Repomix.
The content has been processed where content has been compressed (code blocks are separated by ⋮---- delimiter).

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Content has been compressed - code blocks are separated by ⋮---- delimiter
- Files are sorted by Git change count (files with more changes are at the bottom)

## Additional Info

# Directory Structure
```
cmd/
  main.go
config/
  proxy_pool_config.toml
internal/
  api/
    fetch.go
    http.go
    rate_test.go
    rate.go
  config/
    config_test.go
    config.go
  log/
    log.go
  models/
    domaininfo_test.go
    domaininfo.go
    ipdetails.go
    pool.go
  pool/
    init.go
.gitignore
go.mod
README.md
repomix.config.json
```

# Files

## File: repomix.config.json
````json
{
  "output": {
    "filePath": "repomix-output.md",
    "style": "markdown",
    "parsableStyle": false,
    "fileSummary": true,
    "directoryStructure": true,
    "removeComments": false,
    "removeEmptyLines": false,
    "compress": false,
    "topFilesLength": 5,
    "showLineNumbers": false,
    "copyToClipboard": false,
    "git": {
      "sortByChanges": true,
      "sortByChangesMaxCommits": 100
    }
  },
  "include": [],
  "ignore": {
    "useGitignore": true,
    "useDefaultPatterns": true,
    "customPatterns": []
  },
  "security": {
    "enableSecurityCheck": true
  },
  "tokenCount": {
    "encoding": "o200k_base"
  }
}
````

## File: cmd/main.go
````go
// Author       :loyd
// Date         :2025-03-30 12:41:41
// LastEditors  :loyd
// LastEditTime :2025-03-31 00:05:46
// Description  :程序入口点
//
⋮----
package main
⋮----
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"proxies-api/internal/api"
	"proxies-api/internal/config"
	"proxies-api/internal/log"
	"proxies-api/internal/pool"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)
⋮----
"context"
"fmt"
"net/http"
"os"
"os/signal"
"proxies-api/internal/api"
"proxies-api/internal/config"
"proxies-api/internal/log"
"proxies-api/internal/pool"
"strconv"
"syscall"
"time"
⋮----
"github.com/gin-gonic/gin"
⋮----
// Function: main
⋮----
// Purpose: 程序入口函数，初始化并启动代理IP池系统
⋮----
// Parameters:
//   - 无
⋮----
// Returns: 无
func main()
⋮----
// 初始化配置和IP池
⋮----
// 设置日志系统
⋮----
// 初始化并启动HTTP服务
⋮----
// 启动IP获取服务
⋮----
// 监听系统信号以优雅关闭
⋮----
// Function: setupLogSystem
⋮----
// Purpose: 配置并初始化日志系统
⋮----
func setupLogSystem()
⋮----
// 创建日志目录
⋮----
// 配置日志系统
⋮----
// Function: setupHTTPServer
⋮----
// Purpose: 设置并启动HTTP服务器
⋮----
// Returns: (*gin.Engine, *http.Server) - Gin引擎实例和HTTP服务器实例
func setupHTTPServer() (*gin.Engine, *http.Server)
⋮----
// 设置gin模式
⋮----
// 创建路由器
⋮----
// 注册中间件
⋮----
// 注册路由
⋮----
// 获取配置的端口
⋮----
port = 8000 // 默认端口
⋮----
// 创建HTTP服务器
⋮----
// 启动HTTP服务
⋮----
// Function: startIPFetchService
⋮----
// Purpose: 启动IP获取服务
⋮----
func startIPFetchService()
⋮----
// 延迟一秒等待其他服务初始化完成
⋮----
// 开始更新IP库
⋮----
// Function: setupGracefulShutdown
⋮----
// Purpose: 设置优雅关闭机制
⋮----
//   - router (*gin.Engine): Gin引擎实例
//   - server (*http.Server): HTTP服务器实例
⋮----
func setupGracefulShutdown(router *gin.Engine, server *http.Server)
⋮----
// 创建系统信号通道
⋮----
// 等待信号
⋮----
// 创建一个超时上下文用于优雅关闭
⋮----
// 关闭所有现有连接，但等待请求完成
⋮----
// 停止IP获取服务
⋮----
// 关闭IP池和其他资源
````

## File: config/proxy_pool_config.toml
````toml
[server]
port = 8000

# 代理池配置
[pool]
size = 100

# 自定义代理源
[proxy_source]
urls = [
    "http://example.com/proxy1"
]
request_frequency = 10
local_path = [
    "/var/proxy_pool/proxy.txt",
    "/var/proxy_pool/proxy2.txt",
    "/var/proxy_pool/proxy3.txt",
]
read_frequency = 10

# 单独对每个 ip 设置, delay 也可以使用 URL 参数动态设置
[ip]
lifetime_seconds = 900
delay = [0, 5]
````

## File: internal/api/fetch.go
````go
// Author       :loyd
// Date         :2025-03-30 18:58:57
// LastEditors  :loyd
// LastEditTime :2025-03-30 20:04:22
// Description  :从互联网中或者本地文件夹读取
//
⋮----
package api
⋮----
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"proxies-api/internal/config"
	"proxies-api/internal/pool"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/strutil"
)
⋮----
"fmt"
"io"
"net/http"
"os"
"proxies-api/internal/config"
"proxies-api/internal/pool"
"strconv"
"strings"
"sync"
"time"
⋮----
"github.com/duke-git/lancet/v2/strutil"
⋮----
// IPFetchService 定义IP获取服务结构体
type IPFetchService struct {
	file_ticker   *time.Ticker  // 文件读取定时器
	url_ticker    *time.Ticker  // URL请求定时器
	stop_chan     chan struct{} // 停止信号通道
⋮----
file_ticker   *time.Ticker  // 文件读取定时器
url_ticker    *time.Ticker  // URL请求定时器
stop_chan     chan struct{} // 停止信号通道
ticker_mutex  sync.Mutex    // 互斥锁保护定时器操作
is_running    bool          // 服务是否在运行
paths         []string      // 本地文件路径
urls          []string      // 远程URL列表
read_interval time.Duration // 文件读取间隔
req_interval  time.Duration // URL请求间隔
⋮----
// 全局IP获取服务实例
var ip_fetch_service *IPFetchService
var service_once sync.Once
⋮----
// 重试相关常量
const (
	DEFAULT_MAX_RETRIES    = 3
	DEFAULT_RETRY_INTERVAL = 1 * time.Second
	DEFAULT_MAX_INTERVAL   = 30 * time.Second
	DEFAULT_TIMEOUT        = 30 * time.Second
)
⋮----
// RetryConfig 重试配置结构体
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
	MaxInterval   time.Duration // 最大重试间隔
	Timeout       time.Duration // 请求超时时间
}
⋮----
MaxRetries    int           // 最大重试次数
RetryInterval time.Duration // 重试间隔
MaxInterval   time.Duration // 最大重试间隔
Timeout       time.Duration // 请求超时时间
⋮----
// 默认重试配置
var defaultRetryConfig = RetryConfig{
	MaxRetries:    DEFAULT_MAX_RETRIES,
	RetryInterval: DEFAULT_RETRY_INTERVAL,
	MaxInterval:   DEFAULT_MAX_INTERVAL,
	Timeout:       DEFAULT_TIMEOUT,
}
⋮----
// 默认时间间隔（秒）
const DEFAULT_READ_INTERVAL = 60
const DEFAULT_REQUEST_INTERVAL = 300
⋮----
/**
 * Function: retrieveIPFetchService
 *
 * Purpose: 获取单例IPFetchService实例
 *
 * Parameters:
 *   - 无
 *
 * Returns: *IPFetchService - IP获取服务实例
 */
func retrieveIPFetchService() *IPFetchService
⋮----
// 从配置读取参数
⋮----
// 验证时间间隔，确保大于0
⋮----
// 创建IP获取服务实例
⋮----
/**
 * Function: RefreshIPRegistry
 *
 * Purpose: 启动IP获取服务，从文件和URL中周期性获取代理IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func RefreshIPRegistry()
⋮----
/**
 * Function: StopIPFetchService
 *
 * Purpose: 停止IP获取服务
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止服务
 */
func StopIPFetchService() bool
⋮----
/**
 * Function: Start
 *
 * Purpose: 启动IP获取服务，定期从文件和URL获取IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功启动服务
 */
func (s *IPFetchService) Start() bool
⋮----
// 如果已经在运行，则返回
⋮----
// 再次检查时间间隔是否大于0
⋮----
// 确保IP池已初始化
⋮----
// 立即执行一次搜索
⋮----
// 创建定时器
⋮----
// 启动后台协程执行定时任务
⋮----
/**
 * Function: Stop
 *
 * Purpose: 停止IP获取服务
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止服务
 */
func (s *IPFetchService) Stop() bool
⋮----
// 检查服务是否正在运行
⋮----
// 停止定时器
⋮----
// 发送停止信号
⋮----
// 更新状态
⋮----
/**
 * Function: NeedUpdateIPRegistry
 *
 * Purpose: 检查代理池中的IP数量是否低于阈值，判断是否需要更新IP库
 *
 * Parameters: 无
 *
 * Returns: bool - 如果IP数量低于阈值返回true表示需要更新，否则返回false
 */
func NeedUpdateIPRegistry() bool
⋮----
// 检查池大小与可用IP
⋮----
// 获取配置中设置的最小IP数量要求
⋮----
// 如果当前可用IP数量低于阈值，则需要更新
⋮----
/**
 * Function: searchFromFiles
 *
 * Purpose: 从文件中读取并解析IP地址
 *
 * Parameters:
 *   - paths ([]string): 文件路径列表
 *
 * Returns: 无
 */
func searchFromFiles(paths []string)
⋮----
// 读取配置文件
⋮----
// 如果找到了文件，打印找到的文件路径
⋮----
// 从文件中加载
⋮----
// 从文件读取内容
⋮----
// 解析并处理IP
⋮----
/**
 * Function: ContainsIgnoreCase
 *
 * Purpose: 检查字符串是否包含子字符串，忽略大小写
 *
 * Parameters:
 *   - source (string): 源字符串
 *   - substr (string): 要检查的子字符串
 *
 * Returns: bool - 如果源字符串包含子字符串（忽略大小写）则返回true，否则返回false
 */
func ContainsIgnoreCase(source string, substr string) bool
⋮----
// 转换为小写后进行比较
⋮----
/**
 * Function: searchFromURLs
 *
 * Purpose: 从URL中获取并解析IP地址，支持网页内容解析和错误重试
 *
 * Parameters:
 *   - urls ([]string): URL列表
 *
 * Returns: 无
 */
func searchFromURLs(urls []string)
⋮----
// 创建自定义HTTP客户端，添加超时设置
⋮----
// 创建请求对象
⋮----
// 添加用户代理头部，降低被拒绝的可能性
⋮----
// 使用重试机制发送请求
var resp *http.Response
var retryErr error
⋮----
// 计算退避时间
⋮----
// 发送HTTP请求
⋮----
// 读取响应内容
⋮----
// 根据内容类型进行不同处理
⋮----
// 解析HTML内容中的IP
⋮----
// 解析普通文本内容
⋮----
/**
 * Function: parseHTMLContent
 *
 * Purpose: 解析HTML内容中的IP地址
 *
 * Parameters:
 *   - html_content (string): HTML内容
 *   - ip_regex (string): 用于匹配IP的正则表达式
 *
 * Returns: 无
 */
func parseHTMLContent(html_content string, ip_regex string)
⋮----
// 移除HTML标签，只保留文本内容
⋮----
// 移除多余空白
⋮----
// 解析清理后的文本中的IP
⋮----
/**
 * Function: findValidFilePaths
 *
 * Purpose: 检验并收集所有有效的文件路径
 *
 * Parameters:
 *   - paths ([]string): 待验证的文件路径列表
 *
 * Returns: []string - 所有存在且有效的文件路径列表，如果没有找到则返回空数组
 */
func findValidFilePaths(paths []string) []string
⋮----
var validated_paths []string
⋮----
// 遍历配置中的本地路径
⋮----
// 检查文件是否存在
⋮----
// 📁 文件存在，添加到有效路径列表
⋮----
/**
 * Function: parseAndProcessIPs
 *
 * Purpose: 使用正则表达式解析文本中的IP地址并添加到代理池
 *
 * Parameters:
 *   - content (string): 要解析的文本内容
 *   - regex (string): 用于匹配IP的正则表达式
 *
 * Returns: 无
 */
func parseAndProcessIPs(content string, regex string)
⋮----
// 匹配所有IP和端口
⋮----
// 输出匹配结果并添加到池中
⋮----
// 有端口的情况
⋮----
// 转换端口为整数
⋮----
// 没有端口的情况，使用默认端口80
⋮----
/**
 * Function: GetProxyPoolStatus
 *
 * Purpose: 获取代理池的当前状态信息
 *
 * Parameters: 无
 *
 * Returns: map[string]interface{} - 包含池状态的映射，如池大小、原始IP数量等
 */
func GetProxyPoolStatus() map[string]interface
⋮----
// 获取原始IP数量
⋮----
// 获取域名数量
⋮----
// 构建状态信息
````

## File: internal/api/http.go
````go
// Author       :loyd
// Date         :2025-03-30 18:19:08
// LastEditors  :loyd
// LastEditTime :2025-03-30 18:21:14
// Description  :通过 URL 参数传递数组的 HTTP 服务
//
⋮----
package api
⋮----
import (
	"strconv"
	"strings"

	"proxies-api/internal/models"
	"proxies-api/internal/pool"

	"github.com/gin-gonic/gin"
)
⋮----
"strconv"
"strings"
⋮----
"proxies-api/internal/models"
"proxies-api/internal/pool"
⋮----
"github.com/gin-gonic/gin"
⋮----
type Query struct {
	Domain string `json:"domain"`
	Delay  []int  `json:"delay,omitempty"`
}
⋮----
/**
 * Function: API_Get
 *
 * Purpose: 通过 URL 参数获取代理 IP 信息，包括处理数组形式的 delay 参数
 *
 * Parameters:
 *   - c (*gin.Context): Gin 框架的上下文对象
 *
 * Returns:
 */
func API_Get(c *gin.Context)
⋮----
// 从 URL 获取域名参数
⋮----
// 获取 delay 参数数组，支持两种方式：
// 1. 逗号分隔：?delay=1,5
// 2. 重复参数：?delay=1&delay=5
var delay_array []int
⋮----
// 方式1：通过逗号分隔的方式获取数组
⋮----
// 方式2：通过重复参数的方式获取数组
⋮----
// 处理查询逻辑
⋮----
// 添加新域名
⋮----
// 如果有传递 delay 参数，设置延迟
⋮----
// 复制 origin ips
⋮----
// 返回一个 IP
⋮----
// 返回结果
⋮----
"delay": delay_array, // 返回解析的 delay 数组，方便调试
⋮----
func API_Status(c *gin.Context)
⋮----
/**
 * Function: RegisterUrlRoutes
 *
 * Purpose: 注册使用 URL 参数的 HTTP 路由
 *
 * Parameters:
 *   - router (*gin.Engine): Gin 引擎实例
 *
 * Returns:
 */
func RegisterUrlRoutes(router *gin.Engine)
````

## File: internal/api/rate_test.go
````go
// Author       :loyd
// Date         :2025-03-30 11:24:58
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:28:50
// Description  :测试基于IP的请求限流中间件
⋮----
package api
⋮----
import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)
⋮----
"net/http"
"net/http/httptest"
"sync"
"testing"
"time"
⋮----
"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
"golang.org/x/time/rate"
⋮----
// 测试常量
const (
	TEST_HOST        = "example.com"
	TEST_RATE_LIMIT  = 5
	TEST_BURST_SIZE  = 10
	TEST_MINI_EXPIRY = 1 // 测试用的较短过期时间（分钟）
)
⋮----
TEST_MINI_EXPIRY = 1 // 测试用的较短过期时间（分钟）
⋮----
// Function: TestNewIPRateLimiter
//
// Purpose: 测试限流器初始化功能
⋮----
// Parameters:
//   - t (*testing.T): 测试对象
⋮----
// Returns: 无
func TestNewIPRateLimiter(t *testing.T)
⋮----
// Function: TestGetLimiter
⋮----
// Purpose: 测试获取主机限流器功能
⋮----
func TestGetLimiter(t *testing.T)
⋮----
// 首次获取会创建新的限流器
⋮----
// 再次获取应返回相同的限流器
⋮----
// 检查hosts映射是否包含此主机
⋮----
// Function: TestConcurrentGetLimiter
⋮----
// Purpose: 测试并发获取限流器场景
⋮----
func TestConcurrentGetLimiter(t *testing.T)
⋮----
const num_goroutines = 100
⋮----
var wg sync.WaitGroup
⋮----
// 验证无论多少goroutine，只创建了一个限流器
⋮----
// Function: TestRateLimiting
⋮----
// Purpose: 测试限流器的实际限流功能
⋮----
func TestRateLimiting(t *testing.T)
⋮----
limiter := NewIPRateLimiter(rate.Limit(2), 1) // 每秒2个请求，突发1个
⋮----
// 第一个请求应通过
⋮----
// 连续的请求应被限流
⋮----
// 等待补充令牌
⋮----
// Function: TestCleanupExpiredLimiters
⋮----
// Purpose: 测试清理过期限流器功能
⋮----
func TestCleanupExpiredLimiters(t *testing.T)
⋮----
// 使用较短的过期时间
const quick_expiry = 100 * time.Millisecond
⋮----
// 创建一个自定义的IPRateLimiter用于测试
⋮----
// 添加一个过期的限流器
⋮----
last_access: time.Now().Add(-quick_expiry * 2), // 设置为已过期
⋮----
// 手动触发清理
⋮----
// 验证限流器已被清理
⋮----
// Function: TestRateLimitMiddleware
⋮----
// Purpose: 测试限流中间件功能
⋮----
func TestRateLimitMiddleware(t *testing.T)
⋮----
// 设置测试环境
⋮----
// 测试缺少domain参数的情况
⋮----
// 测试正常请求
⋮----
// 使用高限制，确保不被限流
⋮----
// 测试被限流的请求
⋮----
// 使用低限制，确保被限流
⋮----
// 第一个请求消耗掉唯一的令牌
⋮----
// 下一个请求应被限流
⋮----
// 恢复默认值
⋮----
// Function: TestSetRateLimit
⋮----
// Purpose: 测试动态设置限流参数功能
⋮----
func TestSetRateLimit(t *testing.T)
⋮----
// 保存原始限流器以便恢复
⋮----
// 测试设置新的限流参数
⋮----
// 恢复原始限流器
````

## File: internal/api/rate.go
````go
// Author       :loyd
// Date         :2025-03-30 11:24:58
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:28:50
// Description  :基于IP的请求限流中间件
//
⋮----
package api
⋮----
import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)
⋮----
"net/http"
"sync"
"time"
⋮----
"github.com/gin-gonic/gin"
"golang.org/x/time/rate"
⋮----
// 限流器配置常量
const (
	DEFAULT_RATE_LIMIT = 10 // 每秒请求数限制
	DEFAULT_BURST_SIZE = 20 // 默认令牌桶容量
	CLEANUP_INTERVAL   = 10 // 清理间隔（分钟）
	LIMITER_EXPIRY     = 60 // 限流器过期时间（分钟）
)
⋮----
DEFAULT_RATE_LIMIT = 10 // 每秒请求数限制
DEFAULT_BURST_SIZE = 20 // 默认令牌桶容量
CLEANUP_INTERVAL   = 10 // 清理间隔（分钟）
LIMITER_EXPIRY     = 60 // 限流器过期时间（分钟）
⋮----
// IPRateLimiter 结构体定义
type IPRateLimiter struct {
	hosts      map[string]*hostLimiter // 主机对应的限流器映射
	mu         *sync.RWMutex           // 读写锁，保护并发访问
	rate_limit rate.Limit              // 每秒可以生成的令牌数
	burst_size int                     // 令牌桶的容量
}
⋮----
hosts      map[string]*hostLimiter // 主机对应的限流器映射
mu         *sync.RWMutex           // 读写锁，保护并发访问
rate_limit rate.Limit              // 每秒可以生成的令牌数
burst_size int                     // 令牌桶的容量
⋮----
// 主机限流器，包含过期时间
type hostLimiter struct {
	limiter     *rate.Limiter // 限流器实例
	last_access time.Time     // 最后访问时间
}
⋮----
limiter     *rate.Limiter // 限流器实例
last_access time.Time     // 最后访问时间
⋮----
// NewIPRateLimiter 创建一个新的IP限流器
⋮----
// 功能: 初始化并返回一个新的IPRateLimiter实例
⋮----
// 参数:
//   - rate_limit (rate.Limit): 每秒允许的请求数
//   - burst_size (int): 令牌桶大小
⋮----
// 返回值: *IPRateLimiter - 初始化后的限流器实例
func NewIPRateLimiter(rate_limit rate.Limit, burst_size int) *IPRateLimiter
⋮----
// 启动定期清理过期限流器的协程
⋮----
// 全局默认限流器实例
var RateLimiter = NewIPRateLimiter(DEFAULT_RATE_LIMIT, DEFAULT_BURST_SIZE)
⋮----
// GetLimiter 获取特定主机的限流器
⋮----
// 功能: 返回指定主机的限流器，若不存在则创建新的
⋮----
//   - host (string): 主机标识符，通常是域名或IP地址
⋮----
// 返回值: *rate.Limiter - 对应主机的限流器实例
func (r *IPRateLimiter) GetLimiter(host string) *rate.Limiter
⋮----
// 双重检查，防止并发问题
⋮----
// 更新最后访问时间
⋮----
// cleanupExpiredLimiters 清理过期的限流器
⋮----
// 功能: 定期检查并移除长时间未使用的限流器，防止内存泄漏
⋮----
// 参数: 无
⋮----
// 返回值: 无
func (r *IPRateLimiter) cleanupExpiredLimiters()
⋮----
// 删除超过过期时间未使用的限流器
⋮----
/**
 * Function: RateLimitMiddleware
 *
 * Purpose: 为Gin框架提供基于域名参数的请求限流功能，仅对特定路由生效
 *
 * Parameters:
 *
 * Returns: gin.HandlerFunc - Gin中间件处理函数
 */
func RateLimitMiddleware() gin.HandlerFunc
⋮----
// 获取当前请求路径
⋮----
// 仅对 /proxy 路由进行域名检查和限流
⋮----
"code": 429, // 标准的HTTP状态码
⋮----
// 对所有路由继续处理请求
⋮----
// SetRateLimit 设置全局限流参数
⋮----
// 功能: 动态调整全局默认限流器的参数
⋮----
//   - new_limit (int): 新的每秒请求数限制
//   - new_burst (int): 新的令牌桶容量
⋮----
func SetRateLimit(new_limit int, new_burst int)
````

## File: internal/config/config_test.go
````go
// Author       :loyd
// Date         :2025-03-30 12:00:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 12:00:00
// Description  :配置模块测试
⋮----
package config
⋮----
import (
	"os"
	"path/filepath"
	"testing"
)
⋮----
"os"
"path/filepath"
"testing"
⋮----
/**
 * Function: TestLoadConfig_ValidFile
 *
 * Purpose: 测试加载有效配置文件的功能
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_ValidFile(t *testing.T)
⋮----
// 创建临时目录
⋮----
// 创建有效的配置文件
⋮----
// 测试加载配置
⋮----
// 验证配置内容
⋮----
// 验证全局配置已设置
⋮----
/**
 * Function: TestLoadConfig_FileNotFound
 *
 * Purpose: 测试加载不存在的配置文件时的错误处理
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_FileNotFound(t *testing.T)
⋮----
// 测试加载不存在的文件
⋮----
/**
 * Function: TestLoadConfig_InvalidTOML
 *
 * Purpose: 测试加载格式无效的TOML配置文件时的错误处理
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_InvalidTOML(t *testing.T)
⋮----
// 创建无效的TOML配置文件
⋮----
// 测试加载无效配置
````

## File: internal/config/config.go
````go
// Author       :loyd
// Date         :2025-03-29 01:46:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:37:27
// Description  :读取配置文件
//
⋮----
package config
⋮----
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)
⋮----
"fmt"
"os"
"path/filepath"
⋮----
"github.com/pelletier/go-toml"
⋮----
// 定义应用程序配置结构
type Config struct {
	Server      ServerConfig      `toml:"server"`
	Pool        PoolConfig        `toml:"pool"`
	ProxySource ProxySourceConfig `toml:"proxy_source"`
	Ip          IpConfig          `toml:"ip"`
}
⋮----
// 服务器相关配置
type ServerConfig struct {
	Port int `toml:"port"`
}
⋮----
// 代理池相关配置
type PoolConfig struct {
	Size int `toml:"size"`
}
⋮----
// 代理源相关配置
type ProxySourceConfig struct {
	Url              []string `toml:"urls"`
	RequestFrequency int      `toml:"request_frequency"`
	LocalPath        []string `toml:"local_path"`
	ReadFrequency    int      `toml:"read_frequency"`
}
⋮----
// IP相关配置
type IpConfig struct {
	Lifetime_seconds int   `toml:"lifetime_seconds"`
	Delay            []int `toml:"delay"`
}
⋮----
// 全局配置实例
var APP_CONFIG *Config
⋮----
/**
 * Function: LoadConfig
 *
 * Purpose: 从指定路径读取TOML配置文件并解析到Config结构体
 *
 * Parameters:
 *   - config_path (string): 配置文件的路径
 *
 * Returns: (*Config, error) - 返回配置对象指针和可能的错误
 */
func LoadConfig(config_path string) (*Config, error)
⋮----
// 获取绝对路径
⋮----
// 读取配置文件
⋮----
// 解析TOML配置到结构体
⋮----
// 设置全局配置
⋮----
/**
 * Function: GetConfig
 *
 * Purpose: 获取全局配置实例
 *
 * Parameters: 无
 *
 * Returns: *Config - 返回全局配置对象指针
 */
func GetConfig() *Config
````

## File: internal/log/log.go
````go
// Author       :loyd
// Date         :2025-03-29 02:00:56
// LastEditors  :loyd
// LastEditTime :2025-03-30 19:14:09
// Description  :日志工具封装
//
⋮----
package log
⋮----
import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)
⋮----
"io"
"os"
"path"
"runtime"
"strconv"
⋮----
"github.com/sirupsen/logrus"
"gopkg.in/natefinch/lumberjack.v2"
⋮----
// 常量定义
const (
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)
⋮----
// LogConfig 日志配置结构体
type LogConfig struct {
	LogLevel   string // 日志级别，可选值: debug, info, warn, error
	LogFile    string // 日志文件路径，为空则只输出到控制台
	MaxSize    int    // 单个日志文件最大大小(MB)，超过后会轮转
	MaxBackups int    // 最大备份日志文件数量
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}
⋮----
LogLevel   string // 日志级别，可选值: debug, info, warn, error
LogFile    string // 日志文件路径，为空则只输出到控制台
MaxSize    int    // 单个日志文件最大大小(MB)，超过后会轮转
MaxBackups int    // 最大备份日志文件数量
MaxAge     int    // 日志文件保留天数
Compress   bool   // 是否压缩旧日志文件
⋮----
// 定义私有的全局logger实例
var log = logrus.New()
⋮----
// 初始化日志
func init()
⋮----
// 默认配置 - 在SetupLogger被显式调用前使用
⋮----
// SetupLogger 配置日志系统
func SetupLogger(config LogConfig)
⋮----
// 创建日志目录
⋮----
// 配置日志轮转
⋮----
MaxSize:    config.MaxSize, // 单位MB
⋮----
MaxAge:     config.MaxAge, // 天数
⋮----
// 同时输出到控制台和文件
⋮----
// 设置日志格式
⋮----
// 设置日志级别
⋮----
// 对外暴露的日志方法，直接调用
⋮----
// Debug 输出调试日志
func Debug(format string, args ...interface
⋮----
// Info 输出信息日志
func Info(format string, args ...interface
⋮----
// Warn 输出警告日志
func Warn(format string, args ...interface
⋮----
// Error 输出错误日志
func Error(format string, args ...interface
⋮----
// Fatal 输出致命错误日志
func Fatal(format string, args ...interface
⋮----
// WithFields 创建带字段的日志条目
func WithFields(fields logrus.Fields) *logrus.Entry
````

## File: internal/models/domaininfo_test.go
````go
package models
⋮----
import (
	"fmt"
	"testing"
	"time"
)
⋮----
"fmt"
"testing"
"time"
⋮----
// Function: TestGetIPCount
//
// Purpose: 测试DomainInfo结构体的GetIPCount方法
⋮----
// Parameters:
//   - t (*testing.T): 测试上下文对象
⋮----
// Returns:
func TestGetIPCount(t *testing.T)
⋮----
// 测试用例1: 空IP池
⋮----
// 创建一个空的DomainInfo实例
⋮----
// 验证IP数量为0
⋮----
// 测试用例2: 含有多个IP的池
⋮----
// 创建DomainInfo实例
⋮----
// 添加3个IP
⋮----
// 验证IP数量为3
⋮----
// 删除一个IP后再次验证
⋮----
// 测试用例3: 测试IP重复添加
⋮----
// 多次添加相同IP
⋮----
// OrderedMap会覆盖相同键的值，所以数量应该为1
⋮----
// Function: TestGet
⋮----
// Purpose: 测试DomainInfo结构体的Get方法
⋮----
func TestGet(t *testing.T)
⋮----
// 尝试获取IP
⋮----
// 验证返回错误且IP为nil
⋮----
// 测试用例2: 单个IP
⋮----
// 添加一个IP
old_time := time.Now().Add(-1 * time.Hour) // 1小时前
⋮----
// 记录当前域名的最后使用时间
⋮----
// 获取IP
⋮----
// 验证获取的IP
⋮----
// 验证last_use_time已更新
⋮----
// 验证域名的最后使用时间也已更新
⋮----
// 测试用例3: 多个IP并验证轮换
⋮----
// 添加三个IP
⋮----
domain_info.Add(ip1) // 此时顺序: [ip1]
domain_info.Add(ip2) // 此时顺序: [ip1, ip2]
domain_info.Add(ip3) // 此时顺序: [ip1, ip2, ip3]
⋮----
// 获取第一个IP
⋮----
// 获取第二个IP
⋮----
// 获取第三个IP
⋮----
// 再次获取，应该回到第一个IP（轮换）
⋮----
// Function: TestAdd
⋮----
// Purpose: 测试DomainInfo结构体的Add方法
⋮----
func TestAdd(t *testing.T)
⋮----
// 测试用例1: 向空的DomainInfo添加IP
⋮----
// 初始IP数量应为0
⋮----
// 验证IP已被添加
⋮----
// 验证可以获取到添加的IP
⋮----
// 测试用例2: 添加多个不同的IP
⋮----
// 添加三个不同的IP
⋮----
// 验证IP数量
⋮----
// 获取所有IP并验证是否都被正确添加
⋮----
// 创建一个map来检查所有IP是否存在
⋮----
// 验证所有添加的IP都在集合中
⋮----
// 测试用例3: 重复添加相同IP
⋮----
// 验证IP数量为1
⋮----
// 创建新的IpDetails但使用相同的IP地址
⋮----
// 验证IP数量仍为1（因为使用了相同的键）
⋮----
// 验证存储的是最新添加的IP对象
⋮----
// 测试用例4: 测试添加和删除后的IP位置顺序
⋮----
// 获取所有IP验证添加顺序
⋮----
// 确认IP顺序为 ip1, ip2, ip3
⋮----
// 删除中间的IP (ip2)
⋮----
// 验证IP数量为2
⋮----
// 获取所有IP验证删除后的顺序
⋮----
// 确认IP顺序为 ip1, ip3
⋮----
// 添加新的IP
⋮----
// 获取所有IP验证添加后的顺序
⋮----
// 确认IP顺序为 ip1, ip3, ip4
⋮----
// 验证轮换获取IP的顺序
````

## File: internal/models/domaininfo.go
````go
package models
⋮----
import (
	"fmt"
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)
⋮----
"fmt"
"proxies-api/internal/log"
"sync"
"time"
⋮----
"github.com/duke-git/lancet/v2/maputil"
⋮----
// 常量定义
const (
	ERROR_NO_IP_AVAILABLE = "没有可用的IP"
)
⋮----
// DomainInfo 域名信息结构体
type DomainInfo struct {
	lock           sync.RWMutex                            // 读写锁，用于保护并发访问
	domain         string                                  // 域名
	last_use_time  time.Time                               // 上次使用时间
	delay          []int                                   // 延迟时间
	ip_details_map *maputil.OrderedMap[string, *IpDetails] // IP池，存储与域名相关的IP详情
}
⋮----
lock           sync.RWMutex                            // 读写锁，用于保护并发访问
domain         string                                  // 域名
last_use_time  time.Time                               // 上次使用时间
delay          []int                                   // 延迟时间
ip_details_map *maputil.OrderedMap[string, *IpDetails] // IP池，存储与域名相关的IP详情
⋮----
// Function: NewDomainInfo
//
// Purpose: 创建并初始化一个新的域名信息对象
⋮----
// Parameters:
//   - domain (string): 域名字符串
⋮----
// Returns: *DomainInfo - 初始化完成的DomainInfo实例指针
func NewDomainInfo(domain string) *DomainInfo
⋮----
// Function: SetDelay
⋮----
// Purpose: 设置域名的延迟时间
⋮----
//   - delay (int): 要设置的延迟时间（毫秒）
⋮----
// Returns:
func (d *DomainInfo) SetDelay(delay []int) bool
⋮----
// 超过 2 个延迟时间则不设置
⋮----
// 设置延迟时间
⋮----
// Function: GetDomain
⋮----
// Purpose: 获取DomainInfo对象存储的域名
⋮----
// Returns: string - 域名字符串
func (d *DomainInfo) GetDomain() string
⋮----
// Function: GetLastUseTime
⋮----
// Purpose: 获取域名最后一次被使用的时间
⋮----
// Returns: time.Time - 域名最后使用的时间点
func (d *DomainInfo) GetLastUseTime() time.Time
⋮----
// Function: GetIPCount
⋮----
// Purpose: 获取域名下可用IP的数量
⋮----
// Returns: int - 可用IP的总数量
func (d *DomainInfo) GetIPCount() int
⋮----
// Function: Get
⋮----
// Purpose: 从域名的IP池中获取一个IP并将其移至队列末尾
⋮----
//   - *IpDetails - 获取到的IP详情对象，如果没有可用IP则为nil
//   - error - 可能的错误，如无可用IP时返回错误信息
func (d *DomainInfo) Get() (*IpDetails, error)
⋮----
// 使用读锁进行并发保护
⋮----
// 检查IP池是否为空
⋮----
// 获取第一个元素
⋮----
// 移动到末尾（删除并重新添加）
⋮----
// 更新使用时间
⋮----
// Function: Add
⋮----
// Purpose: 向域名的IP池中添加新的IP
⋮----
//   - ip (*IpDetails): 要添加的IP详情对象指针
⋮----
func (d *DomainInfo) Add(ip *IpDetails)
⋮----
// 使用写锁进行并发保护
⋮----
// 添加IP到有序映射中
⋮----
// Function: Delete
⋮----
// Purpose: 从域名的IP池中删除指定IP
⋮----
//   - ip (*IpDetails): 要删除的IP详情对象指针
⋮----
func (d *DomainInfo) Delete(ip *IpDetails)
⋮----
// 检查IP是否存在并删除
⋮----
// Function: DeleteByIP
⋮----
// Purpose: 通过IP地址字符串从域名的IP池中删除指定IP
⋮----
//   - ipAddress (string): 要删除的IP地址字符串
⋮----
// Returns: bool - 删除操作是否成功，成功返回true，失败返回false
func (d *DomainInfo) DeleteByIP(ipAddress string) bool
⋮----
// Function: GetAllIPs
⋮----
// Purpose: 获取域名下所有的IP详情列表
⋮----
// Returns: []*IpDetails - IP详情对象指针的切片，包含所有IP信息
func (d *DomainInfo) GetAllIPs() []*IpDetails
````

## File: internal/models/ipdetails.go
````go
package models
⋮----
import (
	"fmt"
	"net/http"
	"net/url"
	"proxies-api/internal/log"
	"time"

	"github.com/duke-git/lancet/v2/netutil"
)
⋮----
"fmt"
"net/http"
"net/url"
"proxies-api/internal/log"
"time"
⋮----
"github.com/duke-git/lancet/v2/netutil"
⋮----
// 重试相关常量
const (
	IP_CHECK_MAX_RETRIES    = 2
	IP_CHECK_RETRY_INTERVAL = 1 * time.Second
	IP_CHECK_MAX_INTERVAL   = 5 * time.Second
	IP_CHECK_TIMEOUT        = 5 * time.Second
	IP_CHECK_SUCCESS_CODES  = 200
)
⋮----
// IPCheckConfig IP检查配置结构体
type IPCheckConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
	MaxInterval   time.Duration // 最大重试间隔
	Timeout       time.Duration // 请求超时时间
	SuccessCodes  int           // 成功的状态码
	TestURL       string        // 测试URL
}
⋮----
MaxRetries    int           // 最大重试次数
RetryInterval time.Duration // 重试间隔
MaxInterval   time.Duration // 最大重试间隔
Timeout       time.Duration // 请求超时时间
SuccessCodes  int           // 成功的状态码
TestURL       string        // 测试URL
⋮----
// 默认IP检查配置
var defaultIPCheckConfig = IPCheckConfig{
	MaxRetries:    IP_CHECK_MAX_RETRIES,
	RetryInterval: IP_CHECK_RETRY_INTERVAL,
	MaxInterval:   IP_CHECK_MAX_INTERVAL,
	Timeout:       IP_CHECK_TIMEOUT,
	SuccessCodes:  IP_CHECK_SUCCESS_CODES,
	TestURL:       "https://www.baidu.com",
}
⋮----
// IpDetails IP详情结构体
type IpDetails struct {
	ip            string    // IP 地址
	port          int       // 端口号
	acquire_time  time.Time // 获取时间
	last_use_time time.Time // 上次使用时间
}
⋮----
ip            string    // IP 地址
port          int       // 端口号
acquire_time  time.Time // 获取时间
last_use_time time.Time // 上次使用时间
⋮----
// Function: GetIP
//
// Purpose: 获取IP详情对象中存储的IP地址
⋮----
// Parameters:
⋮----
// Returns: string - IP地址字符串
func (i *IpDetails) GetIP() string
⋮----
func (i *IpDetails) CheckConnection() bool
⋮----
// 创建请求配置
⋮----
// 创建客户端配置
⋮----
// 创建HTTP客户端
⋮----
// 重试逻辑
var lastErr error
⋮----
// 计算退避时间
⋮----
// 发送请求
⋮----
// 检查状态码
⋮----
// 成功，关闭响应体并返回
⋮----
// 所有重试都失败
⋮----
// Function: GetPort
⋮----
// Purpose: 获取IP详情对象中存储的端口号
⋮----
// Returns: int - 端口号整数值
func (i *IpDetails) GetPort() int
⋮----
// Function: GetAcquireTime
⋮----
// Purpose: 获取IP被获取的时间点
⋮----
// Returns: time.Time - IP首次获取的时间
func (i *IpDetails) GetAcquireTime() time.Time
⋮----
// Function: GetLastUseTime
⋮----
// Purpose: 获取IP最后一次被使用的时间点
⋮----
// Returns: time.Time - IP最后一次使用的时间
func (i *IpDetails) GetLastUseTime() time.Time
⋮----
// Function: NewIpDetails
⋮----
// Purpose: 创建并初始化一个新的IP详情对象
⋮----
//   - ip (string): IP地址字符串
//   - port (int): IP的端口号
⋮----
// Returns: *IpDetails - 初始化完成的IpDetails实例指针
func NewIpDetails(ip string, port int) *IpDetails
````

## File: internal/models/pool.go
````go
package models
⋮----
import (
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)
⋮----
"proxies-api/internal/log"
"sync"
"time"
⋮----
"github.com/duke-git/lancet/v2/maputil"
⋮----
// Pool 结构体定义了代理IP池的核心数据结构
type Pool struct {
	size           int64                                       // 池的大小
	ip_ttl         int64                                       // 每个IP的生命周期 seconds
	origin_ip_map  *maputil.ConcurrentMap[string, *IpDetails]  // 原始IP映射
	domains        *maputil.ConcurrentMap[string, *DomainInfo] // 域名信息映射
	clean_ticker   *time.Ticker                                // 清理IP的定时器
	ticker_stopped bool                                        // 标记定时器是否已停止
	ticker_mutex   sync.Mutex                                  // 保护定时器操作的互斥锁
	stop_chan      chan struct{}                               // 用于停止定时器的通道
⋮----
size           int64                                       // 池的大小
ip_ttl         int64                                       // 每个IP的生命周期 seconds
origin_ip_map  *maputil.ConcurrentMap[string, *IpDetails]  // 原始IP映射
domains        *maputil.ConcurrentMap[string, *DomainInfo] // 域名信息映射
clean_ticker   *time.Ticker                                // 清理IP的定时器
ticker_stopped bool                                        // 标记定时器是否已停止
ticker_mutex   sync.Mutex                                  // 保护定时器操作的互斥锁
stop_chan      chan struct{}                               // 用于停止定时器的通道
⋮----
/**
 * Function: GetOriginIPs
 *
 * Purpose: 获取池中所有原始IP及其详情的映射
 *
 * Parameters:
 *
 * Returns: map[string]*IpDetails - IP地址到IP详情对象的映射，包含池中所有原始IP信息
 */
func (p *Pool) GetOriginIPs() map[string]*IpDetails
⋮----
/**
 * Function: CleanTicker
 *
 * Purpose: 启动一个定时器定期清理池中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) CleanTicker()
⋮----
// 如果已经有定时器在运行，则返回
⋮----
// 创建定时器和停止通道
⋮----
// 在后台启动清理任务
⋮----
// 执行清理操作
⋮----
// 收到停止信号，结束goroutine
⋮----
/**
 * Function: cleanExpiredIPs
 *
 * Purpose: 清理所有域名中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) cleanExpiredIPs()
⋮----
// 遍历所有域名
⋮----
// 检查域名的IP池是否为空
⋮----
// 清理过期IP
⋮----
/**
 * Function: StopCleanTicker
 *
 * Purpose: 停止IP清理定时器
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止定时器
 */
func (p *Pool) StopCleanTicker() bool
⋮----
// 检查定时器是否已存在且未停止
⋮----
// 停止定时器
⋮----
// 发送停止信号给goroutine
⋮----
// 标记为已停止
⋮----
/**
 * Function: AddIPToPool
 *
 * Purpose: 向池中添加新的IP和端口，并分配到所有域名
 *
 * Parameters:
 *   - ip (string): 要添加的IP地址
 *   - port (int): 要添加的端口号
 *
 * Returns: 无
 */
func (p *Pool) AddIPToPool(ip string, port int)
⋮----
// 添加原始IP到池中
⋮----
// 添加到全部 domains
⋮----
/**
 * Function: DeleteIPFromPool
 *
 * Purpose: 从池中删除指定IP，并从所有域名中移除该IP
 *
 * Parameters:
 *   - ip (string): 要删除的IP地址
 *
 * Returns:
 */
func (p *Pool) DeleteIPFromPool(ip string)
⋮----
// 删除原始IP从池中
⋮----
// 从所有域名中删除
⋮----
/**
 * Function: GetSize
 *
 * Purpose: 获取IP池的大小
 *
 * Parameters:
 *
 * Returns: int64 - IP池的大小，表示可以容纳的最大域名数量
 */
func (p *Pool) GetSize() int64
⋮----
/**
 * Function: GetDomains
 *
 * Purpose: 获取IP池中所有域名及其信息的映射
 *
 * Parameters:
 *
 * Returns: map[string]*DomainInfo - 域名到域名信息对象的映射
 */
func (p *Pool) GetDomains() map[string]*DomainInfo
⋮----
/**
 * Function: GetDomainInfo
 *
 * Purpose: 获取指定域名的详细信息
 *
 * Parameters:
 *   - domain (string): 要查询的域名
 *
 * Returns:
 *   - *DomainInfo - 域名信息对象，如不存在则为nil
 *   - bool - 表示域名是否存在于池中
 */
func (p *Pool) GetDomainInfo(domain string) (*DomainInfo, bool)
⋮----
/**
 * Function: AddDomain
 *
 * Purpose: 向IP池中添加新的域名
 *
 * Parameters:
 *   - domain (string): 要添加的域名
 *
 * Returns:
 */
func (p *Pool) AddDomain(domain string)
⋮----
// 检查域名是否已存在
⋮----
// 创建新的域名信息对象并添加到池中
⋮----
/**
 * Function: DeleteDomain
 *
 * Purpose: 从IP池中删除指定域名及其所有IP信息
 *
 * Parameters:
 *   - domain (string): 要删除的域名
 *
 * Returns:
 */
func (p *Pool) DeleteDomain(domain string)
⋮----
// 检查域名是否存在
⋮----
// 删除域名信息对象
⋮----
/**
 * Function: NewPool
 *
 * Purpose: 创建并初始化一个新的池对象
 *
 * Parameters:
 *   - size (int64): 池的大小
 *   - lifetime_per_ip (int64): 每个IP的生命周期(秒)
 *   - urlResource (string): 资源地址字符串
 *
 * Returns: *Pool - 初始化完成的Pool实例指针
 */
func NewPool(size int64, lifetime_per_ip int64) *Pool
````

## File: internal/pool/init.go
````go
// Author       :loyd
// Date         :2025-03-29 13:35:01
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:59:44
// Description  :内存池启动
//
⋮----
package pool
⋮----
import (
	"fmt"
	"os"
	"proxies-api/internal/config"
	"proxies-api/internal/models"
	"sync"
)
⋮----
"fmt"
"os"
"proxies-api/internal/config"
"proxies-api/internal/models"
"sync"
⋮----
// 全局内存池实例
var pool_instance *models.Pool
⋮----
// 默认配置文件路径列表
var DEFAULT_CONFIG_PATH = []string{
	"../config/proxy_pool_config.toml",
	"../../config/proxy_pool_config.toml",
	"./config/proxy_pool_config.toml",
}
var once sync.Once
⋮----
// Function: InitializeConfigAndPool
⋮----
// Purpose: 初始化配置文件并创建内存池实例，确保只初始化一次
⋮----
// Parameters:
//   - 无
⋮----
// Returns: void - 无返回值
func InitializeConfigAndPool() *models.Pool
⋮----
// 确保只初始化一次
⋮----
// 配置文件
⋮----
// 如果没有找到配置文件，报错并退出
⋮----
// 初始化内存池
⋮----
// 启动定时器清理过期IP
⋮----
// Function: AccessPool
⋮----
// Purpose: 获取内存池实例，如果未初始化则先初始化
⋮----
// Returns: *models.Pool - 返回内存池实例指针
func AccessPool() *models.Pool
````

## File: .gitignore
````
# ---------- 通用配置 ----------

# 日志

# 测试文件
**/test/

# 临时文件
**/tmp/
**/*.tmp

# 二进制文件
**/*.o
**/*.pyc
**/*.swp
**/*.class
**/*.exe
**/*.dll
**/*.so
**/*.dylib
**/*.a

# 压缩文件
**/*.rar
**/*.tgz
**/*.zip
**/*.7z
**/*.gz
**/*.bz
**/*.xz
**/*.tar
**/*.tar.gz

# 输出目录
**/dist/
**/build/
**/out/
**/target/

# 特殊文件
**/.DS_Store
**/.idea/
**/Thumb.db

# vscode 相关配置
**/.vscode/

# 数据目录
**/data*/
**/archive*/

# ---------- Python 相关 ----------
# 环境
**/.venv/

# python 缓存
**/__pycache__/
**/.pytest_cache/

# ---------- 前端相关 ----------
# node_modules
**/node_modules/

# ---------- 自定义 ----------
**/.vscode/
````

## File: go.mod
````
module proxies-api

go 1.23.3

require (
	github.com/pelletier/go-toml v1.9.5
	github.com/sirupsen/logrus v1.9.2
)

require (
	github.com/duke-git/lancet/v2 v2.3.5
	github.com/stretchr/testify v1.8.3
	golang.org/x/time v0.11.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20221208152030-732eee02a75a // indirect
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
````

## File: README.md
````markdown
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

### 获取系统状态

**请求**:
- 路径: `/status`
- 方法: `GET`
- 参数: 无

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "✅ 获取状态成功",
  "status": {
    "domains_count": 3,
    "origins_count": 50,
    "domains": [
      {
        "name": "example.com",
        "ips_count": 50
      },
      {
        "name": "test.com",
        "ips_count": 50
      }
    ]
  }
}
```

### 限流机制

系统内置请求限流机制，防止单个域名请求过多，默认限制为每秒10个请求，突发上限为20个请求。
注意：限流机制仅应用于 `/proxy` 端点，`/status` 端点不受限流影响。

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
│   └── main.go     # 主程序入口
├── config/         # 配置文件
├── internal/       # 内部实现
│   ├── api/        # HTTP接口
│   ├── config/     # 配置解析
│   ├── log/        # 日志工具
│   ├── models/     # 数据结构
│   └── pool/       # IP池管理
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
````
