// Author       :loyd
// Date         :2025-03-30 12:00:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 12:00:00
// Description  :配置模块测试

package config

import (
	"os"
	"path/filepath"
	"testing"
)

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
func TestLoadConfig_ValidFile(t *testing.T) {
	// 创建临时目录
	temp_dir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("❌ 无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(temp_dir)

	// 创建有效的配置文件
	config_content := `
[server]
port = 8080

[pool]
size = 100

[proxy]
url = "http://example.com/proxy"

[ip]
lifetime_seconds = 3600
delay = [10, 20, 30]
`
	config_path := filepath.Join(temp_dir, "config.toml")
	err = os.WriteFile(config_path, []byte(config_content), 0644)
	if err != nil {
		t.Fatalf("❌ 无法写入测试配置文件: %v", err)
	}

	// 测试加载配置
	config, err := LoadConfig(config_path)
	if err != nil {
		t.Fatalf("❌ 加载有效配置文件应成功，但返回错误: %v", err)
	}

	// 验证配置内容
	if config.Server.Port != 8080 {
		t.Errorf("❌ 服务器端口配置错误，期望 8080，实际为 %d", config.Server.Port)
	}
	if config.Pool.Size != 100 {
		t.Errorf("❌ 代理池大小配置错误，期望 100，实际为 %d", config.Pool.Size)
	}
	if config.ProxySource.Url != "http://example.com/proxy" {
		t.Errorf("❌ 代理URL配置错误，期望 'http://example.com/proxy'，实际为 '%s'",
			config.ProxySource.Url)
	}
	if config.Ip.Lifetime_seconds != 3600 {
		t.Errorf("❌ IP生命周期配置错误，期望 3600，实际为 %d",
			config.Ip.Lifetime_seconds)
	}
	if len(config.Ip.Delay) != 3 || config.Ip.Delay[0] != 10 ||
		config.Ip.Delay[1] != 20 || config.Ip.Delay[2] != 30 {
		t.Errorf("❌ 延迟配置错误，期望 [10,20,30]，实际为 %v", config.Ip.Delay)
	}

	// 验证全局配置已设置
	global_config := GetConfig()
	if global_config != config {
		t.Error("❌ 全局配置未正确设置")
	}
}

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
func TestLoadConfig_FileNotFound(t *testing.T) {
	// 测试加载不存在的文件
	_, err := LoadConfig("/path/to/nonexistent/config.toml")
	if err == nil {
		t.Error("❌ 加载不存在的配置文件应返回错误，但没有")
	}
}

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
func TestLoadConfig_InvalidTOML(t *testing.T) {
	// 创建临时目录
	temp_dir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("❌ 无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(temp_dir)

	// 创建无效的TOML配置文件
	invalid_content := `
[server]
port = "not an integer" # 应该是整数，不是字符串
`
	invalid_path := filepath.Join(temp_dir, "invalid.toml")
	err = os.WriteFile(invalid_path, []byte(invalid_content), 0644)
	if err != nil {
		t.Fatalf("❌ 无法写入测试配置文件: %v", err)
	}

	// 测试加载无效配置
	_, err = LoadConfig(invalid_path)
	if err == nil {
		t.Error("❌ 加载格式无效的配置文件应返回错误，但没有")
	}
}
