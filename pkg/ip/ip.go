package ip

import (
	"fmt"
	"net"
	"supreme-flamego/pkg/colorful"
)

func GetLocalHost() (res []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(colorful.Red("net.Interfaces failed, err: " + err.Error()))
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						res = append(res, ipnet.IP.String())
					}
				}
			}
		}

	}
	return
}
