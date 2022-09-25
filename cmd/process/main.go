package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gbrlbrbs/CES27_Lab1/internal/types"
	"github.com/gbrlbrbs/CES27_Lab1/internal/utils"
)

// global variables
var safeLogicalClock types.SafeLogicalClock
var insideCS bool
var waiting bool
var receivedAllReplies bool
var repliesReceived []int

var queuedRequests []int

// function that verifies if a specific process is in the queue
func searchInList(processID int, repliesReceived []int) bool {
	for _, i := range repliesReceived {
		if i == processID {
			return true
		}
	}
	return false
}

//function that verifies between two process which one has priority to access CS
func amIPriority(processID, processClock, thisID, thisClock int) bool {

	if thisClock < processClock {
		return true
	} else if thisClock > processClock {
		return false
	} else {
		if thisID < processID {
			return true
		} else {
			return false
		}
	}
}

func readInputFromStdin(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, err := reader.ReadLine()
		utils.PrintError(err)
		if err == nil {
			ch <- string(text)
		}
	}
}

// function that creates the queue os processes that have request CS while CS was occupied
// It queues request from requesting process without replying
func queueRequest(processID int, requestsInLine []int) (newQueue []int) {
	newQueue = append(requestsInLine, processID)
	return
	//auxiliar := append(requestsInLine, processID)
	//fmt.Println("auxiliar:", len(auxiliar))
	//return auxiliar
	//return append(requestsInLine, processID)
}

// function to build and send reply messages
func replyToProcess(thisClock, thisID, processID int, connections []*net.UDPConn) {
	// convert to string
	str_clock := strconv.Itoa(thisClock)
	str_id := strconv.Itoa(thisID)

	// utils.Concatenate
	message, buf := utils.Concatenate(str_id, str_clock, "reply")

	// send reply to process
	index := processID - 1
	_, err := connections[index].Write(buf)
	utils.PrintErrorAndMessage(err, message)
}

// function to send to all the other processes a request to use CS
func askOtherProcessToUseCS(thisClock, thisID int, connections []*net.UDPConn) {

	if thisID == 2 {
		time.Sleep(time.Second * 2)
	}

	str_clock := strconv.Itoa(thisClock)
	str_id := strconv.Itoa(thisID)

	// utils.Concatenate
	message, buf := utils.Concatenate(str_id, str_clock, "request")

	//Multicast request to all N-1 processes
	for _, conn2process := range connections {
		if conn2process != nil {
			_, err := conn2process.Write(buf)
			utils.PrintErrorAndMessage(err, message)
		}
	}
}

// function to send message to CS and sleep (all other processes have replied)
func sendMessageToCS(thisClock, thisID int, text string, sharedResource *net.UDPConn) {
	insideCS = true
	str_clock := strconv.Itoa(thisClock)
	str_id := strconv.Itoa(thisID)

	// utils.Concatenate
	message, buf := utils.Concatenate(str_id, str_clock, text)

	//send message to shared resource
	_, err := sharedResource.Write(buf)
	utils.PrintErrorAndMessage(err, message)
	//sleep
	time.Sleep(time.Second * 15)
}

func replyQueuedRequests(thisClock, thisID int, connections []*net.UDPConn) {
	str_clock := strconv.Itoa(thisClock)
	str_id := strconv.Itoa(thisID)

	// utils.Concatenate
	message, buf := utils.Concatenate(str_id, str_clock, "reply")

	//Reply to all queued processes
	for _, request_id := range queuedRequests {
		index := request_id - 1
		//reply queued request
		_, err := connections[index].Write(buf)
		if err != nil {
			fmt.Println(">>>")
			fmt.Println(message, err)
		}
	}

	queuedRequests = nil
}

// function to do the actions of a process leaving CS: change state to "released" and reset the flags
func releaseCS(thisClock, thisID int, connections []*net.UDPConn) {
	insideCS = false
	waiting = false
	receivedAllReplies = false

	//reply to any queued request
	replyQueuedRequests(thisClock, thisID, connections)

	//clear reply received list
	repliesReceived = nil
}

func RicartAgrawala(thisClock, thisID int, text string, connections []*net.UDPConn, sharedResource *net.UDPConn) {

	waiting = true
	askOtherProcessToUseCS(thisClock, thisID, connections)

	//Wait until received N-1 replies
	fmt.Println(">>>I am waiting for all N-1 replies")
	for !receivedAllReplies {
	}

	// Enter CS after receive all N-1 replies
	fmt.Println(">>>Got inside CS")

	// Send your message to CS
	sendMessageToCS(thisClock, thisID, text, sharedResource)
	fmt.Println(">>>Just sent my message to CS")

	// Leave CS
	releaseCS(thisClock, thisID, connections)
	fmt.Println(">>>Got out of CS!")
}

func readFromUDP(thisID, nServers int, serverConn *net.UDPConn, connections []*net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, _, err := serverConn.ReadFromUDP(buf)
		utils.PrintError(err)
		// aux == "<id>,<clock>,<type>"
		messageStr := string(buf[:n])
		// messageArr == ["<id>", "<clock>", "<type>"]
		messageArr := strings.Split(messageStr, ",")
		processIDStr := messageArr[0]
		processClockStr := messageArr[1]
		messageType := messageArr[2]
		fmt.Println(">>>Received message from process", processIDStr, "having logical clock =", processClockStr, "of type =", messageType)
		processID, _ := strconv.Atoi(processIDStr)
		processClock, _ := strconv.Atoi(processClockStr)
		thisClock := safeLogicalClock.Get()
		if messageType == "request" {
			if insideCS || (waiting && amIPriority(processID, processClock, thisID, thisClock)) {
				// using or waiting and has the priority
				fmt.Println(">>>Queued request from process", processIDStr)
				queuedRequests = queueRequest(processID, queuedRequests)
			} else {
				// reply immediately
				fmt.Println(">>>Replying process", processIDStr)
				replyToProcess(thisClock, thisID, processID, connections)
			}
		} else if messageType == "reply" {
			if !searchInList(processID, repliesReceived) {
				// append to replies received if not found
				repliesReceived = append(repliesReceived, processID)
			}
			if len(repliesReceived) >= nServers {
				fmt.Println(">>>Received all replies")
				receivedAllReplies = true
			}
		} else {
			fmt.Println(">>>Unknown message type; message =", messageStr, "type =", messageType)
		}
		safeLogicalClock.Set(utils.MaxNumber(safeLogicalClock.Get(), processClock) + 1)
		fmt.Println(">>>Logical clock now set to", safeLogicalClock.Get())
	}
}

func processListenUDP(procPort string) (serverConn *net.UDPConn) {
	serverAddr, err := net.ResolveUDPAddr("udp", procPort)
	utils.PrintError(err)

	serverConn, err = net.ListenUDP("udp", serverAddr)
	utils.PrintError(err)
	return
}

func makeConnections(processID int, nServers int, ports []string) (connections []*net.UDPConn, sharedResource *net.UDPConn) {
	connections = make([]*net.UDPConn, nServers+1)
	for i := 0; i <= nServers; i++ {
		if i == processID-1 {
			connections[i] = nil
			continue
		}

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
	utils.PrintError(err)
	return
}

func getArgs() (id int, procPort string, ports []string, nServers int) {
	id, err := strconv.Atoi(os.Args[1])
	utils.CheckError(err)
	procPort = os.Args[id+1]
	ports = os.Args[2:]
	nServers = len(os.Args) - 3
	return
}

func main() {
	id, procPort, ports, nServers := getArgs()
	connections, sharedResource := makeConnections(id, nServers, ports)
	serverConn := processListenUDP(procPort)
	safeLogicalClock = types.SafeLogicalClock{} // will create with logicalClock: 0

	ch := make(chan string)
	go readInputFromStdin(ch)
	for {
		go readFromUDP(id, nServers, serverConn, connections)
		select {
		case x, valid := <-ch:
			if valid {
				compare, _ := strconv.Atoi(x)
				if compare != id && x == "x" {
					if insideCS || waiting {
						fmt.Println(">>>Ignored")
					} else {
						thisClock := safeLogicalClock.Get()
						fmt.Println(">>>Requesting access with id =", id, "and logical clock =", thisClock)
						go RicartAgrawala(thisClock, id, "sugou", connections, sharedResource)
					}
				} else {
					safeLogicalClock.Increase()
					fmt.Println(">>>Logical clock increased to", safeLogicalClock.Get())
				}
			} else {
				fmt.Println(">>>Channel closed")
			}
		default:
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Second * 1)
	}
}
