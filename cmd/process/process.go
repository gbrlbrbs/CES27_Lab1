package main

import (
	"net"
	"os"
	"strconv"

	"github.com/gbrlbrbs/CES27_Lab1/internal/utils"
)

func makeConnections(nServers int, ports []string) (connections []*net.UDPConn, sharedResource *net.UDPConn) {
	connections = make([]*net.UDPConn, nServers)
	for i := 0; i < nServers; i++ {
		port := ports[i]

		serverAddr, err := net.ResolveUDPAddr("udp", port)
		utils.PrintError(err)

		localAddr, err := net.ResolveUDPAddr("udp", ":0")
		utils.PrintError(err)

		connections[i], err = net.DialUDP("udp", localAddr, serverAddr)
		utils.PrintError(err)
	}

	csAddr, err := net.ResolveUDPAddr("udp", ":10001")
	utils.PrintError(err)

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	utils.PrintError(err)

	sharedResource, err = net.DialUDP("udp", localAddr, csAddr)
	return
}

func getArgs() (id int, procPort string, ports []string, nServers int) {
	id, err := strconv.Atoi(os.Args[0])
	utils.CheckError(err)
	procPort = os.Args[id]
	ports = os.Args[1:]
	nServers = len(os.Args) - 1
	return
}

func main() {
	id, procPort, ports, nServers := getArgs()
	connections, sharedResource := makeConnections(nServers, ports)
}