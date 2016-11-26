package main

import "net"
import "fmt"
import "runtime"
import "log"

import "os"
import "syscall"
import "os/signal"

import "msgs"

// import "config"

type Handler struct {
	server_sock net.Listener
	max_workers int
	max_buff    int
	msgChannel  chan net.Conn
	killsig     chan bool
	dispatcher  *msgs.MsgDispatcher
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func NewHandler() *Handler {
	port := ":1337"
	workers := runtime.NumCPU() * 10
	fmt.Println("Num workers: ", workers)
	buff_size := 1000

	// listen on all interfaces
	msgChannel := make(chan net.Conn, buff_size)
	listenerSock, _ := net.Listen("tcp", port)
	dispatcher := msgs.NewMsgDispatcher(msgChannel, workers)

	return &Handler{
		server_sock: listenerSock,
		max_workers: workers,
		max_buff:    buff_size,
		msgChannel:  msgChannel,
		killsig:     make(chan bool),
		dispatcher:  dispatcher}
}

func (handler *Handler) Run() {
	fmt.Println("Starting Handler...")
	handler.dispatcher.Run()

	handler.signalHandler()

	handler.serve()
	<-handler.killsig
	fmt.Println("Handler closed")
}

func (handler *Handler) serve() {
	chanCap := 0
	for {
		conn, err := handler.server_sock.Accept()
		if err != nil {
			switch errType := err.(type) {
			case *net.OpError:
				if errType.Op == "accept" {
					println("Server socket closed")
					return
				}

			default:
				fmt.Println(err)
			}
		}

		handler.msgChannel <- conn

		chanCap = getMax(chanCap, len(handler.msgChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func (handler *Handler) Close() {
	fmt.Println("Closing handler")
	handler.dispatcher.Close()

	handler.server_sock.Close()
	handler.killsig <- true
}

func (handler *Handler) signalHandler() {
	killsig := make(chan os.Signal, 1)
	signal.Notify(killsig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-killsig
		handler.Close()
		os.Exit(1)
	}()
}

func getMax(num1, num2 int) int {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func setupLogger() {

	direxists, err := exists("./logs")
	if !direxists {
		direrr := os.Mkdir("./logs", 0766)
		if direrr != nil {
			fmt.Printf("Error init logging dir: %v", direrr)
			os.Exit(1)
		}
	}

	fexists, err := exists("./logs/handler_log.txt")
	if !fexists {
		os.Create("./logs/handler_log.txt")
	}

	errlog, err := os.Open("./logs/handler_log.txt")
	if err != nil {
		fmt.Printf("Error initalizing error logging: %v", err)
		os.Exit(1)
	}

	Trace = log.New(errlog, "Application Log: ", log.Lshortfile|log.LstdFlags)

}

// *** MAIN *** //
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setupLogger()
	handler := NewHandler()
	handler.Run()
}
