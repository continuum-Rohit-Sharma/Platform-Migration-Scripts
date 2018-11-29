package rest

import "net"

// GetNetworkInterfaces returns list of all tcp interfaces of application
func GetNetworkInterfaces(url string) ([]string, error) {
	host, port, err := net.SplitHostPort(url)
	if err != nil {
		return nil, err
	}

	// If host specified directly, no reason to continue
	if host != "" {
		return []string{host + ":" + port}, nil
	}

	ipSlice, err := getIPs()
	if err != nil {
		return nil, err
	}

	networkInterface := make([]string, 0, len(ipSlice))
	for _, ip := range ipSlice {
		networkInterface = append(networkInterface, ip+":"+port)
	}

	return networkInterface, nil
}

// getIPs returns list of all available IPs (v4 only)
func getIPs() ([]string, error) {
	// In most cases we will have 3 and more interfaces
	IPs := make([]string, 0, 3)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				IPs = append(IPs, ip.String())
			}
		}
	}
	return IPs, nil
}
