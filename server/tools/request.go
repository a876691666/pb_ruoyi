package tools

import (
	"net"
	"net/http"
)

// GetIPAddr 获取请求的真实IP地址（不依赖任何 header）
// 实现说明：直接使用 r.RemoteAddr 并用 net.SplitHostPort 解析出 host 部分。
// 这在不依赖于代理或 headers 的场景下能取得连接对端的 IP 地址。
func GetIPAddr(r *http.Request) string {
	if r == nil {
		return ""
	}

	// r.RemoteAddr 通常形如 "IP:port" 或 "[IPv6]:port"。
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	// 如果解析失败，直接返回 RemoteAddr（可能不含端口）
	return r.RemoteAddr
}

// GetLocationByIP 根据 IP 地址获取地理位置（示例函数，需自行实现具体逻辑）
func GetLocationByIP(ip string) string {
	// TODO: 实现根据 IP 获取地理位置的逻辑
	return ""
}
