package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"bufio"
	"strings"
	"github.com/gbrlbrbs/CES27_Lab1/internal/utils"

)



// function that verifies if a specific process is in the queue
func searchInList(process_id int) bool{
	for _,i := range replies_received {
		if i == process_id{
			return true
		}
	}
	return false
}

//function that verifies between two process which one has priority to acess CS
func amIthepriority(process_id int, process_clock int) bool {

	if request_clock < process_clock {//my logical clock at request is lower than another's
		return true//I am priority
	} else if process_clock > request_clock {//my logical clock at request isn't lower than another's
		return false//I'm not priority 
	} else {//my logical clock at request is equal to another's
		if request_id < process_id {//my id is lower
			return true//I am priority
		} else {//my id isn't lower
			return false//I'm not priority
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

//function that creates the queue os processes that have request CS qhile CS was occupied
//It queues request from requesting process without replying
func queueRequest(process_id int){
	queued_request = append(queued_request,process_id)
}

//function to build and send reply messages
func reply2process(process_id int){
	str_clock:= strconv.Itoa(mylogicalClock)//convert clock (int) into a string type
	str_id := strconv.Itoa(id)//convert id (int) into a string type
	
	// concatenate
	message = concatenate(str_id , str_clock, "reply")
	
	// send reply to process
	index := process_id - 1
     _,err := CliConn[index].Write(buf)
     if err != nil {
        fmt.Println(message, err)
	}
}


// function to send to all the other processes a request to use CS
func askOtherProcessToUseCS(clock int){

	str_clock := strconv.Itoa(clock)//convert clock (int) into a string type 
	str_id := strconv.Itoa(id)//convert id (int) into a string type
	
	// concatenate
	message = concatenate(str_id , str_clock, "request")

	//Multicast request to all N-1 processes
	for _, conn2process := range CliConn {

     	_,err := conn2process.Write(buf)
     	if err != nil {
        	fmt.Println(message, err)
		}
	}
}

//function to send message to CS and sleep (all other processes have replied)
func sendMessageToCS(request_clock int, text string){
	insideCS = true
	str_request_clock := strconv.Itoa(lc_requisicao)//convert clock (int) into a string type 
	str_id := strconv.Itoa(id)//convert id (int) into a string type
	
	// concatenate
	message = concatenate(str_id , str_clock, text)

	//send message to shared resource
     _,err := sharedResource.Write(buf)
     if err != nil {
        fmt.Println(message, err)
	}
	//sleep
	time.Sleep(time.Second*3)
}


func replyQueuedRequests(){
	str_clock := strconv.Itoa(mylogicalClock)//convert clock (int) into a string type 
	str_id := strconv.Itoa(id)//convert id (int) into a string type
	
	// concatenate
	message = concatenate(str_id , str_clock, "reply")

	//Reply to all queued processes
	for _,id:= range queued_request {
		index := id -1
		//reply queued request
		_,err := CliConn[index].Write(buf)
		if err != nil {
		   fmt.Println(message, err)
	   }
   }
}

// function to do the actions of a process leaving CS: change state to "released" and reset the flags
func releaseCS(){
	insideCS = false
	waiting = false
	received_all_replies = false

	//reply to any queued request
	replyQueuedRequests()

	//clear reply received list
	replies_received = nil
}

func Ricart_Agrawala(lc_requisicao int, text_simples string){

	waiting = true
	askOtherProcessToUseCS(request_clock)
	
	//Wait until received N-1 replies
	fmt.Println("I am waiting for all N-1 replies")
	for !received_all_replies {}

	// Enter CS after receive all N-1 replies
	fmt.Println("Got inside CS")

	// Send your message to CS
	sendMessageToCS(request_clock,text)
	fmt.Println("Just sent my message to CS")

	// Leave CS
	releaseCS()
	fmt.Println("Got out of CS!")
}

func processListenUDP(procPort string) (serverConn *net.UDPConn) {
	serverAddr, err := net.ResolveUDPAddr("udp", procPort)
	utils.PrintError(err)

	serverConn, err = net.ListenUDP("udp", serverAddr)
	utils.PrintError(err)
	return
}

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
	serverConn := processListenUDP(procPort)
	logicalClock := 0

	// close connections on main function finish
	defer serverConn.Close()
	for i := nServers; i < nServers; i++ {
		defer connections[i].Close()
	}

	ch := make(chan string)
	go readInputFromStdin(ch)
	for {
		
	}
}