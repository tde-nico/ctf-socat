package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func handleConnection(conn net.Conn, cmd string, f string) {
	args := []string{cmd}
	cmdObj := exec.Command(args[0], args[1:]...)
	cmdObj.Env = append(os.Environ(), "FLAG="+f)

	stdin, err := cmdObj.StdinPipe()
	if err != nil {
		log.Fatalf("Error creating stdin pipe: %v\n", err)
		return
	}

	stdout, err := cmdObj.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe: %v\n", err)
		return
	}

	if err := cmdObj.Start(); err != nil {
		log.Fatalf("Error executing command: %v\n", err)
		return
	}

	go func() {
		io.Copy(stdin, conn)
		stdin.Close()
	}()
	io.Copy(conn, stdout)

	cmdObj.Wait()
	conn.Close()
}

func serve(port string, cmd string, f string) {
	init_regex(f)
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Printf("Listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v\n", err)
			continue
		}
		debugLog("%v\n", conn.RemoteAddr())

		unique_flag := getUniqueFlag()
		go handleConnection(conn, cmd, unique_flag)
	}
}

func main() {
	var port string
	var cmd string
	var f string

	flag.StringVar(&port, "p", "1337", "Specify the port")
	flag.StringVar(&cmd, "c", "", "Specify the command")
	flag.StringVar(&f, "f", "flag{test_[random 6]}", "Specify the flag")
	flag.BoolVar(&debug, "d", false, "Debug mode")
	flag.Parse()

	if cmd == "" {
		flag.Usage()
		os.Exit(1)
	}

	serve(port, cmd, f)
}
