// Author       :loyd
// Date         :2025-03-29 13:35:01
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:59:44
// Description  :内存池启动
//

package pool

import (
	"fmt"
	"os"
	"proxies-api/internal/config"
	"proxies-api/internal/models"
	"sync"
)

// 全局内存池实例
var pool_instance *models.Pool

// 默认配置文件路径列表
var DEFAULT_CONFIG_PATH = []string{
	"../config/proxy_pool_config.toml",
	"../../config/proxy_pool_config.toml",
	"./config/proxy_pool_config.toml",
}
var once sync.Once

// Function: InitializeConfigAndPool
//
// Purpose: 初始化配置文件并创建内存池实例，确保只初始化一次
//
// Parameters:
//   - 无
//
// Returns: void - 无返回值
func InitializeConfigAndPool() *models.Pool {
	// 确保只初始化一次
	once.Do(func() {
		// 配置文件
		found := false
		for _, path := range DEFAULT_CONFIG_PATH {
			if _, err := os.Stat(path); err == nil {
				config.LoadConfig(path)
				found = true
				fmt.Println("🔍 找到配置文件:", path)
				break
			}
		}

		// 如果没有找到配置文件，报错并退出
		if !found {
			fmt.Println("❌ 错误: 未找到配置文件，请确保以下路径至少有一个有效的配置文件")
			for _, path := range DEFAULT_CONFIG_PATH {
				fmt.Println("  - " + path)
			}
			os.Exit(1)
		}

		// 初始化内存池
		pool_instance = models.NewPool(
			int64(config.APP_CONFIG.Pool.Size),
			int64(config.APP_CONFIG.Ip.Lifetime_seconds))

		fmt.Println("✅ 内存池初始化成功，容量:", config.APP_CONFIG.Pool.Size)

		// 启动定时器清理过期IP
		pool_instance.CleanTicker()
		fmt.Println("⏰ IP清理定时器已启动")
	})

	return pool_instance
}

// Function: AccessPool
//
// Purpose: 获取内存池实例，如果未初始化则先初始化
//
// Parameters:
//   - 无
//
// Returns: *models.Pool - 返回内存池实例指针
func AccessPool() *models.Pool {
	InitializeConfigAndPool()
	return pool_instance
}
