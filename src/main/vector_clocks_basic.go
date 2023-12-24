package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	serverPort = ":5553"
	balance    = 1000
)

var (
	vectorTime = []int{0, 0, 0}
	serverList = []string{"server1", "server2", "server3"}
	serverMap  = map[string]int{"server1": 0, "server2": 1, "server3": 2}
	mutex      sync.Mutex
)

func main() {
	go startServer()
	time.Sleep(3 * time.Second)
	go startClient()

	// Prevent the main function from exiting
	select {}
}

func startServer() {
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var msg string
	fmt.Fscan(conn, &msg)

	data := strings.Split(msg, "|")
	sendingServer := data[0]
	currentServer := data[1]
	amount, _ := strconv.Atoi(data[2])
	vectorIncoming := strings.Split(data[3], ",")

	mutex.Lock()
	defer mutex.Unlock()
	index := serverMap[currentServer]
	vectorTime[index]++

	for i, v := range vectorIncoming {
		incomingTime, _ := strconv.Atoi(v)
		vectorTime[i] = max(vectorTime[i], incomingTime)
	}

	fmt.Printf("Received %d from %s. Vector time: %v\n", amount, sendingServer, vectorTime)
}

func startClient() {
	for {
		time.Sleep(5 * time.Second)
		selectEvent()
	}
}

func selectEvent() {
	switch rand.Intn(3) {
	case 0:
		fmt.Println("Deposit Selected")
		deposit()
	case 1:
		fmt.Println("Withdraw Selected")
		withdraw()
	case 2:
		fmt.Println("Transfer Selected")
		transfer()
	}
}

func deposit() {
	amount := rand.Intn(100)
	fmt.Printf("Depositing %d\n", amount)

	mutex.Lock()
	defer mutex.Unlock()
	updateVectorTime()
	fmt.Printf("Balance after deposit: %d. Vector time: %v\n", balance+amount, vectorTime)
}

func withdraw() {
	amount := rand.Intn(100)
	fmt.Printf("Withdrawing %d\n", amount)

	mutex.Lock()
	defer mutex.Unlock()
	if balance-amount > 0 {
		updateVectorTime()
		fmt.Printf("Balance after withdrawal: %d. Vector time: %v\n", balance-amount, vectorTime)
	} else {
		fmt.Println("Insufficient funds, cannot withdraw")
	}
}

func transfer() {
	amount := rand.Intn(100)
	transferTo := rand.Intn(3)

	mutex.Lock()
	updateVectorTime()
	mutex.Unlock()

	fmt.Printf("Transferring %d to server %d\n", amount, transferTo)
	msg := fmt.Sprintf("%s|%s|%d|%s", "currentServer", serverList[transferTo], amount, vectorTimeToString())

	conn, err := net.Dial("tcp", serverList[transferTo]+serverPort)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Fprint(conn, msg)
}

func updateVectorTime() {
	// Update the vector time for the current server
	vectorTime[serverMap["currentServer"]]++
}

func vectorTimeToString() string {
	vectorString := make([]string, len(vectorTime))
	for i, v := range vectorTime {
		vectorString[i] = strconv.Itoa(v)
	}
	return strings.Join(vectorString, ",")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
