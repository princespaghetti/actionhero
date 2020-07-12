package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// ActionMessage is a struct representing a message from CSM
type ActionMessage struct {
	Version             int    `json:"Version"`
	ClientID            string `json:"ClientId"`
	Type                string `json:"Type"`
	Service             string `json:"Service"`
	Action              string `json:"Api"`
	Timestamp           int    `json:"Timestamp"`
	AttemptLatency      int    `json:"AttemptLatency"`
	Fqdn                string `json:"Fqdn"`
	UserAgent           string `json:"UserAgent"`
	AccessKey           string `json:"AccessKey"`
	Region              string `json:"Region"`
	HTTPStatusCode      int    `json:"HttpStatusCode"`
	FinalHTTPStatusCode int    `json:"FinalHttpStatusCode"`
	XAmzRequestID       string `json:"XAmzRequestId"`
	XAmzID2             string `json:"XAmzId2"`
}

func listen(connection *net.UDPConn, quit chan struct{}) {
	buffer := make([]byte, 1024)
	n, _, err := 0, new(net.UDPAddr), error(nil)
	var message ActionMessage
	for err == nil {
		n, _, err = connection.ReadFromUDP(buffer)
		err := json.Unmarshal(buffer[:n], &message)
		if err != nil {
			log.Println(err)
		}
		//Each action taken sends two json messages. The first has a type of "ApiCallAttempt" this filters for the API call itself
		if message.Type == "ApiCall" {
			fmt.Println(strings.ToLower(message.Service) + ":" + message.Action)
		}
	}
	fmt.Println("listener failed - ", err)
	quit <- struct{}{}
}

//SetupCloseHandler Displays a message when the user closes the program
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCtrl+C pressed, Stopping...")
		os.Exit(0)
	}()
}

func main() {

	var port = 31000
	var err error
	if os.Getenv("AWS_CSM_PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("AWS_CSM_PORT"))
		if err != nil {
			fmt.Println("Could not parse value of AWS_CSM_PORT Exiting...")
			os.Exit(1)
		}
	}
	addr := net.UDPAddr{
		Port: port,
		IP:   net.IP{127, 0, 0, 1},
	}
	connection, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Could not start Action hero on the specified port, Exiting...")
		os.Exit(1)
	}
	fmt.Println("Action Hero Starting...")
	SetupCloseHandler()
	quit := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go listen(connection, quit)
	}
	<-quit // hang until an error
}
