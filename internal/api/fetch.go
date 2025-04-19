// Author       :loyd
// Date         :2025-03-30 18:58:57
// LastEditors  :loyd
// LastEditTime :2025-03-30 20:04:22
// Description  :从互联网中或者本地文件夹读取
//

package api

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

// IPFetchService 定义IP获取服务结构体
type IPFetchService struct {
	file_ticker   *time.Ticker  // 文件读取定时器
	url_ticker    *time.Ticker  // URL请求定时器
	stop_chan     chan struct{} // 停止信号通道
	ticker_mutex  sync.Mutex    // 互斥锁保护定时器操作
	is_running    bool          // 服务是否在运行
	paths         []string      // 本地文件路径
	urls          []string      // 远程URL列表
	read_interval time.Duration // 文件读取间隔
	req_interval  time.Duration // URL请求间隔
}

// 全局IP获取服务实例
var ip_fetch_service *IPFetchService
var service_once sync.Once

// 重试相关常量
const (
	DEFAULT_MAX_RETRIES    = 3
	DEFAULT_RETRY_INTERVAL = 1 * time.Second
	DEFAULT_MAX_INTERVAL   = 30 * time.Second
	DEFAULT_TIMEOUT        = 30 * time.Second
)

// RetryConfig 重试配置结构体
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
	MaxInterval   time.Duration // 最大重试间隔
	Timeout       time.Duration // 请求超时时间
}

// 默认重试配置
var defaultRetryConfig = RetryConfig{
	MaxRetries:    DEFAULT_MAX_RETRIES,
	RetryInterval: DEFAULT_RETRY_INTERVAL,
	MaxInterval:   DEFAULT_MAX_INTERVAL,
	Timeout:       DEFAULT_TIMEOUT,
}

// 默认时间间隔（秒）
const DEFAULT_READ_INTERVAL = 60
const DEFAULT_REQUEST_INTERVAL = 300

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
func retrieveIPFetchService() *IPFetchService {
	service_once.Do(func() {
		// 从配置读取参数
		req_frequency := config.GetConfig().ProxySource.RequestFrequency
		read_frequency := config.GetConfig().ProxySource.ReadFrequency
		paths := config.GetConfig().ProxySource.LocalPath
		urls := config.GetConfig().ProxySource.Url

		// 验证时间间隔，确保大于0
		if req_frequency <= 0 {
			fmt.Printf("⚠️ 请求频率配置无效: %d，使用默认值: %d秒\n",
				req_frequency, DEFAULT_REQUEST_INTERVAL)
			req_frequency = DEFAULT_REQUEST_INTERVAL
		}

		if read_frequency <= 0 {
			fmt.Printf("⚠️ 读取频率配置无效: %d，使用默认值: %d秒\n",
				read_frequency, DEFAULT_READ_INTERVAL)
			read_frequency = DEFAULT_READ_INTERVAL
		}

		// 创建IP获取服务实例
		ip_fetch_service = &IPFetchService{
			stop_chan:     make(chan struct{}),
			is_running:    false,
			paths:         paths,
			urls:          urls,
			read_interval: time.Duration(read_frequency) * time.Second,
			req_interval:  time.Duration(req_frequency) * time.Second,
		}
	})

	return ip_fetch_service
}

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
func RefreshIPRegistry() {
	service := retrieveIPFetchService()
	service.Start()
}

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
func StopIPFetchService() bool {
	service := retrieveIPFetchService()
	return service.Stop()
}

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
func (s *IPFetchService) Start() bool {
	s.ticker_mutex.Lock()
	defer s.ticker_mutex.Unlock()

	// 如果已经在运行，则返回
	if s.is_running {
		fmt.Println("⚠️ IP获取服务已经在运行中")
		return false
	}

	// 再次检查时间间隔是否大于0
	if s.read_interval <= 0 {
		fmt.Println("⚠️ 读取间隔必须大于0，使用默认间隔60秒")
		s.read_interval = DEFAULT_READ_INTERVAL * time.Second
	}

	if s.req_interval <= 0 {
		fmt.Println("⚠️ 请求间隔必须大于0，使用默认间隔300秒")
		s.req_interval = DEFAULT_REQUEST_INTERVAL * time.Second
	}

	// 确保IP池已初始化
	pool.AccessPool()
	fmt.Println("🚀 启动IP获取服务...")

	// 创建新的通道和定时器
	s.stop_chan = make(chan struct{})
	s.file_ticker = time.NewTicker(s.read_interval)
	s.url_ticker = time.NewTicker(s.req_interval)
	s.is_running = true

	// 立即执行一次搜索
	if err := searchFromFiles(s.paths); err != nil {
		fmt.Printf("❌ 文件搜索失败: %v\n", err)
	}
	if err := searchFromURLs(s.urls); err != nil {
		fmt.Printf("❌ URL搜索失败: %v\n", err)
	}

	// 启动后台协程执行定时任务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("❌ IP获取服务发生panic: %v\n", r)
				s.Stop() // 确保服务被正确停止
			}
		}()

		for {
			select {
			case <-s.file_ticker.C:
				fmt.Println("🕒 执行定时文件读取...")
				needUpdate := NeedUpdateIPRegistry()
				if needUpdate {
					if err := searchFromFiles(s.paths); err != nil {
						fmt.Printf("❌ 定时文件读取失败: %v\n", err)
					}
				}
			case <-s.url_ticker.C:
				fmt.Println("🌐 执行定时URL请求...")
				needUpdate := NeedUpdateIPRegistry()
				if needUpdate {
					if err := searchFromURLs(s.urls); err != nil {
						fmt.Printf("❌ 定时URL请求失败: %v\n", err)
					}
				}
			case <-s.stop_chan:
				fmt.Println("🛑 IP获取服务接收到停止信号")
				return
			}
		}
	}()

	fmt.Println("✅ IP获取服务已启动")
	return true
}

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
func (s *IPFetchService) Stop() bool {
	s.ticker_mutex.Lock()
	defer s.ticker_mutex.Unlock()

	// 检查服务是否正在运行
	if !s.is_running {
		fmt.Println("⚠️ IP获取服务未运行")
		return false
	}

	// 停止定时器
	if s.file_ticker != nil {
		s.file_ticker.Stop()
		s.file_ticker = nil
	}

	if s.url_ticker != nil {
		s.url_ticker.Stop()
		s.url_ticker = nil
	}

	// 发送停止信号
	if s.stop_chan != nil {
		select {
		case <-s.stop_chan:
			// 通道已经关闭
		default:
			close(s.stop_chan)
		}
	}

	// 更新状态
	s.is_running = false
	fmt.Println("✅ IP获取服务已停止")

	return true
}

/**
 * Function: NeedUpdateIPRegistry
 *
 * Purpose: 检查代理池中的IP数量是否低于阈值，判断是否需要更新IP库
 *
 * Parameters: 无
 *
 * Returns: bool - 如果IP数量低于阈值返回true表示需要更新，否则返回false
 */
func NeedUpdateIPRegistry() bool {

	// 检查池大小与可用IP
	available_ips := pool.AccessPool().GetOriginIPs()
	available_count := len(available_ips)

	// 获取配置中设置的最小IP数量要求
	MIN_IP_COUNT := int(pool.AccessPool().GetSize())

	// 如果当前可用IP数量低于阈值，则需要更新
	if available_count < MIN_IP_COUNT {
		fmt.Printf("⚠️ 当前可用IP数量(%d)低于设定阈值(%d)，需要更新IP库\n",
			available_count, MIN_IP_COUNT)
		return true
	}

	fmt.Printf("✅ 当前可用IP数量充足: %d/%d\n", available_count, MIN_IP_COUNT)
	return false
}

/**
 * Function: searchFromFiles
 *
 * Purpose: 从文件中读取并解析IP地址
 *
 * Parameters:
 *   - paths ([]string): 文件路径列表
 *
 * Returns: error - 如果搜索失败则返回错误，否则返回nil
 */
func searchFromFiles(paths []string) error {
	// 读取配置文件
	file_paths := findValidFilePaths(paths)
	if len(file_paths) == 0 {
		return fmt.Errorf("❌ 未找到有效的配置文件")
	}

	// 如果找到了文件，打印找到的文件路径
	fmt.Println("✅ 找到配置文件:", file_paths)

	// 从文件中加载
	ip_regex := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::(\d{1,5}))?`

	var lastErr error
	for _, file_path := range file_paths {
		// 从文件读取内容
		file_content, err := os.ReadFile(file_path)
		if err != nil {
			lastErr = fmt.Errorf("❌ 读取文件 %s 失败: %w", file_path, err)
			fmt.Println(lastErr)
			continue
		}

		// 解析并处理IP
		parseAndProcessIPs(string(file_content), ip_regex)
	}

	return lastErr
}

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
func ContainsIgnoreCase(source string, substr string) bool {
	// 转换为小写后进行比较
	return strings.Contains(
		strings.ToLower(source),
		strings.ToLower(substr),
	)
}

/**
 * Function: searchFromURLs
 *
 * Purpose: 从URL中获取并解析IP地址，支持网页内容解析和错误重试
 *
 * Parameters:
 *   - urls ([]string): URL列表
 *
 * Returns: error - 如果搜索失败则返回错误，否则返回nil
 */
func searchFromURLs(urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("❌ 没有配置URL来源")
	}

	ip_regex := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::(\d{1,5}))?`

	// 创建自定义HTTP客户端，添加超时设置
	client := &http.Client{
		Timeout: defaultRetryConfig.Timeout,
	}

	var lastErr error
	for _, url := range urls {
		fmt.Println("🌐 正在从URL获取IP:", url)

		// 创建请求对象
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("❌ 创建请求对象失败: %w", err)
			fmt.Println(lastErr)
			continue
		}

		// 添加用户代理头部，降低被拒绝的可能性
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		// 使用重试机制发送请求
		var resp *http.Response
		var retryErr error
		for retry := 0; retry <= defaultRetryConfig.MaxRetries; retry++ {
			if retry > 0 {
				// 计算退避时间
				backoff := defaultRetryConfig.RetryInterval * time.Duration(1<<uint(retry-1))
				if backoff > defaultRetryConfig.MaxInterval {
					backoff = defaultRetryConfig.MaxInterval
				}
				fmt.Printf("🔄 第 %d 次重试，等待 %v 后重试...\n", retry, backoff)
				time.Sleep(backoff)
			}

			// 发送HTTP请求
			resp, retryErr = client.Do(req)
			if retryErr == nil && resp.StatusCode == http.StatusOK {
				break
			}

			if retryErr != nil {
				fmt.Printf("❌ 第 %d 次请求失败: %v\n", retry+1, retryErr)
			} else {
				fmt.Printf("❌ 第 %d 次请求返回非200状态码: %d %s\n", retry+1, resp.StatusCode, resp.Status)
				resp.Body.Close()
			}

			if retry == defaultRetryConfig.MaxRetries {
				lastErr = fmt.Errorf("❌ 达到最大重试次数 %d，放弃请求", defaultRetryConfig.MaxRetries)
				fmt.Println(lastErr)
				continue
			}
		}

		if retryErr != nil {
			lastErr = fmt.Errorf("❌ 所有重试都失败，放弃请求: %w", retryErr)
			fmt.Println(lastErr)
			continue
		}

		defer resp.Body.Close()

		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("❌ 读取响应失败: %w", err)
			fmt.Println(lastErr)
			continue
		}

		content_type := resp.Header.Get("Content-Type")
		fmt.Printf("📄 内容类型: %s\n", content_type)

		// 根据内容类型进行不同处理
		if ContainsIgnoreCase(content_type, "html") {
			fmt.Println("🔍 检测到HTML内容，进行特殊解析...")
			// 解析HTML内容中的IP
			parseHTMLContent(string(body), ip_regex)
		} else {
			// 解析普通文本内容
			parseAndProcessIPs(string(body), ip_regex)
		}
	}

	return lastErr
}

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
func parseHTMLContent(html_content string, ip_regex string) {
	// 移除HTML标签，只保留文本内容
	text_content := strutil.ReplaceWithMap(html_content, map[string]string{
		`<[^>]*>`: " ",
	})

	// 移除多余空白
	text_content = strutil.ReplaceWithMap(text_content, map[string]string{
		`\s+`: " ",
	})

	fmt.Println("🧹 已清理HTML标签，正在解析IP...")

	// 解析清理后的文本中的IP
	parseAndProcessIPs(text_content, ip_regex)
}

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
func findValidFilePaths(paths []string) []string {
	var validated_paths []string

	// 遍历配置中的本地路径
	for _, path := range paths {
		// 检查文件是否存在
		if _, err := os.Stat(path); err == nil {
			// 📁 文件存在，添加到有效路径列表
			validated_paths = append(validated_paths, path)
		}
	}

	return validated_paths
}

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
func parseAndProcessIPs(content string, regex string) {
	// 确保IP池已初始化
	poolInstance := pool.AccessPool()

	// 匹配所有IP和端口
	matches := strutil.RegexMatchAllGroups(regex, content)

	// 输出匹配结果并添加到池中
	for _, match := range matches {
		if len(match) > 2 && match[2] != "" {
			// 有端口的情况
			ip := match[1]
			port_str := match[2]

			// 转换端口为整数
			port, err := strconv.Atoi(port_str)
			if err != nil {
				fmt.Printf("❌ 端口格式错误: %s:%s\n", ip, port_str)
				continue
			}

			fmt.Printf("✅ 找到IP和端口: %s:%d 并添加到代理池\n", ip, port)
			poolInstance.AddIPToPool(ip, port)
		} else if len(match) > 1 {
			// 没有端口的情况，使用默认端口80
			ip := match[1]
			fmt.Printf("✅ 找到IP: %s (使用默认端口:80) 并添加到代理池\n", ip)
			poolInstance.AddIPToPool(ip, 80)
		}
	}
}

/**
 * Function: GetProxyPoolStatus
 *
 * Purpose: 获取代理池的当前状态信息
 *
 * Parameters: 无
 *
 * Returns: map[string]interface{} - 包含池状态的映射，如池大小、原始IP数量等
 */
func GetProxyPoolStatus() map[string]interface{} {
	poolInstance := pool.AccessPool()

	// 获取原始IP数量
	originIPs := poolInstance.GetOriginIPs()
	originIPCount := len(originIPs)

	// 获取域名数量
	domains := poolInstance.GetDomains()
	domainCount := len(domains)

	// 构建状态信息
	status := map[string]interface{}{
		"pool_size":       poolInstance.GetSize(),
		"origin_ip_count": originIPCount,
		"domain_count":    domainCount,
		"domains":         domains,
	}

	return status
}
