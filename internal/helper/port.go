package helper

import (
	"fmt"
	"net"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort(port int) (int, error) {
	if port > 0 && IsPortAvailable(port) {
		return port, nil
	}

	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil

}

// 判断端口是否可用（未被占用）
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}

	defer listener.Close()
	return true
}

func GetUpstreamPort() (port int, err error) {
	return GetUpstreamPortFromRegistry()
	//pid, err := FindPidByName("SeewoCore.exe")
	//if err != nil {
	//	return 0, err
	//}
	//var outBytes bytes.Buffer
	//cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr LISTENING | findstr %d", pid)
	//cmd := exec.Command("cmd", "/c", cmdStr)
	//cmd.Stdout = &outBytes
	//err = cmd.Run()
	//if err != nil {
	//	return 0, err
	//}
	//resStr := outBytes.String()
	//part := strings.Fields(resStr)
	//if len(part) < 2 {
	//	return 0, fmt.Errorf("no such process")
	//}
	//
	//if len(part[1]) != 0 {
	//	part2 := strings.Split(part[1], ":")
	//	port, err := strconv.Atoi(part2[len(part2)-1])
	//	if err != nil {
	//		return 0, err
	//	}
	//	return port, nil
	//}
	//
	//return 0, fmt.Errorf("unknown error")
}
