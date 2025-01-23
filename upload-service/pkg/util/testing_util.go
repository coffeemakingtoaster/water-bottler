package util

import "net"

func GetAvailablePort() int {
	// Listen on 0 for kernel to assign us a port :)
	ln, err := net.Listen("tcp", ":0")

	defer ln.Close()

	if err != nil {
		panic(err)
	}

	return ln.Addr().(*net.TCPAddr).Port
}
