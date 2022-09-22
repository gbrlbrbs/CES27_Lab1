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

//Variáveis globais interessantes para o processo
var err string
var myPort string //porta do meu servidor
var nServers int //qtde de outros processo
var CliConn []*net.UDPConn //vetor com conexões para os servidores
 //dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo
 //mensagens dos outros processos)

var request_id int //numero identificador do processo
var mylogicalClock int
var insideCS bool
var waiting bool
var received_all_replies bool
var sharedResource *net.UDPConn
var queued_request []int
var request_clock int
var replies_received []int

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
	} 
	else if process_clock > request_clock {//my logical clock at request isn't lower than another's
		return false//I'm not priority 
	} 
	else {//my logical clock at request is equal to another's
		if request_id < process_id {//my id is lower
			return true//I am priority
		} 
		else {//my id isn't lower
			return false//I'm not priority
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
	str_id := strconv.Itoa(request_id)//convert id (int) into a string type
	
	// utils.Concatenate
	message := utils.Concatenate(str_id , str_clock, "reply")
	
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
	str_id := strconv.Itoa(request_id)//convert id (int) into a string type
	
	// utils.Concatenate
	message := utils.Concatenate(str_id , str_clock, "request")

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
	str_clock := strconv.Itoa(request_clock)//convert clock (int) into a string type 
	str_id := strconv.Itoa(request_id)//convert id (int) into a string type
	
	// utils.Concatenate
	message := utils.Concatenate(str_id , str_clock, text)

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
	str_id := strconv.Itoa(request_id)//convert id (int) into a string type
	
	// utils.Concatenate
	message := utils.Concatenate(str_id , str_clock, "reply")

	//Reply to all queued processes
	for _,request_id:= range queued_request {
		index := request_id -1
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

func Ricart_Agrawala(request_clock int, text string){

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

func main(){
	initConnections()
	insideCS = false
	waiting = false
	//O fechamento de conexões deve ficar aqui, assim só fecha
	//conexão quando a main morrer

	defer ServConn.Close()
	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}

	//Todo Process fará a mesma coisa: ouvir msg e mandar infinitos
	//i’s para os outros processos
	ch := make(chan string)
	go readInput(ch)
	for {
			
		//Server
		go doServerJob()

		// When there is a request (from stdin). Do it!
		select {
			case x, valid := <- ch :
				if valid {
                        compare,_ := strconv.Atoi(x)
						if ( compare != id && x == "x"){

							//Ver se esta na CS ou esperando
							if (insideCS || waiting){
								fmt.Println("x ignorado!")
							} else {
								fmt.Printf("Solicitando acesso com ID = %d e Logical Clock = %d\n", request_id, mylogicalClock)
								text := "CS sugou"
								request_clock = mylogicalClock
								go Ricart_Agrawala(request_clock,text)
							}
							
						} else{

							mylogicalClock++
							fmt.Printf("Atualizado logicalClock para %d \n",mylogicalClock)
						}
				} else {
						 
					fmt.Println("Channel closed!")
						
				}
				
			default:
			
				// Do nothing in the non-blocking approach.
			
				time.Sleep(time.Second * 1)
		}
			
		// Wait a while
		time.Sleep(time.Second * 1)
	}
}