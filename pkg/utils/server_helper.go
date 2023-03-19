package utils

import (
	"net"
	"strings"
)

// 获取服务器Ip
// func GetServerIp() (ip string) {

// 	addrs, err := net.InterfaceAddrs()

// 	if err != nil {
// 		return ""
// 	}

// 	for _, address := range addrs {
// 		// 检查ip地址判断是否回环地址
// 		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
// 			if ipNet.IP.To4() != nil {
// 				ip = ipNet.IP.String()
// 			}
// 		}
// 	}

// 	return
// }

/**** 问题：我在本地多网卡机器上，运行分布式场景，此函数返回的ip有误导致rpc连接失败。 遂google结果如下：
 ***	1、https://www.jianshu.com/p/301aabc06972
 ***	2、https://www.cnblogs.com/chaselogs/p/11301940.html
****/

// GetServerNodeIP 获取服务节点ip
func GetServerNodeIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip
	// 当物理机通过该方式获取的ip不正确,故使用上述方法，可能多网卡时使用下面的方法
	// ip, err := externalIP()
	// if err != nil {
	// 	return ""
	// }
	// return ip.String()
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, err
}

func getIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
