// Author       :loyd
// Date         :2025-03-29 01:46:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:37:27
// Description  :读取配置文件
//

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

// 定义应用程序配置结构
type Config struct {
	Server      ServerConfig      `toml:"server"`
	Pool        PoolConfig        `toml:"pool"`
	ProxySource ProxySourceConfig `toml:"proxy_source"`
	Ip          IpConfig          `toml:"ip"`
}

// 服务器相关配置
type ServerConfig struct {
	Port int `toml:"port"`
}

// 代理池相关配置
type PoolConfig struct {
	Size int `toml:"size"`
}

// 代理源相关配置
type ProxySourceConfig struct {
	Url              []string `toml:"urls"`
	RequestFrequency int      `toml:"request_frequency"`
	LocalPath        []string `toml:"local_path"`
	ReadFrequency    int      `toml:"read_frequency"`
}

// IP相关配置
type IpConfig struct {
	Lifetime_seconds int   `toml:"lifetime_seconds"`
	Delay            []int `toml:"delay"`
}

// 全局配置实例
var APP_CONFIG *Config

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
func LoadConfig(config_path string) (*Config, error) {
	// 获取绝对路径
	abs_path, err := filepath.Abs(config_path)
	if err != nil {
		return nil, fmt.Errorf("❌ 无法获取配置文件绝对路径: %w", err)
	}

	// 读取配置文件
	fmt.Printf("📖 正在读取配置文件: %s\n", abs_path)
	config_data, err := os.ReadFile(abs_path)
	if err != nil {
		return nil, fmt.Errorf("❌ 读取配置文件失败: %w", err)
	}

	// 解析TOML配置到结构体
	config := &Config{}
	err = toml.Unmarshal(config_data, config)
	if err != nil {
		return nil, fmt.Errorf("❌ 解析TOML配置失败: %w", err)
	}

	// 设置全局配置
	APP_CONFIG = config

	fmt.Println("✅ 配置加载成功!")
	return config, nil
}

/**
 * Function: GetConfig
 *
 * Purpose: 获取全局配置实例
 *
 * Parameters: 无
 *
 * Returns: *Config - 返回全局配置对象指针
 */
func GetConfig() *Config {
	return APP_CONFIG
}
