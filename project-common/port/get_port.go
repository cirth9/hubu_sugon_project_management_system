package port

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

func GetFreePort() (int, error) {
	// 动态获取可用端口
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	//如果address参数中的端口号为空或为“0”，则自动选择端口号
	l, err := net.Listen("tcp", addr.String())
	if err != nil {
		return 0, err
	}
	zap.S().Info(l.Addr().String())
	return l.Addr().(*net.TCPAddr).Port, nil
}

func FinAvailablePort(startPort, endPort int) (int, error) {
	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			time.Sleep(time.Second)
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available port found in the specified range")
}
