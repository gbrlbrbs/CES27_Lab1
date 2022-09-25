package main

import (
	"fmt"
	"net"
	"strings"

	//"time"

	"github.com/gbrlbrbs/CES27_Lab1/internal/utils"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":10001")
	utils.CheckError(err)
	connection, err := net.ListenUDP("udp", addr)
	utils.CheckError(err)
	defer connection.Close()
	fmt.Println(">>>server listening at", connection.LocalAddr().String())
	for {
		message := make([]byte, 1024)
		rlen, remote, err := connection.ReadFromUDP(message)
		utils.CheckError(err)
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Println(">>>received: msg =", data, "sender =", remote)
	}
}
