package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

type config struct {
	port int
	host string
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	err = runCmd(os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("ccdns", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.IntVar(&c.port, "p", 53, "Port where thet server will listen for incoming requests")
	fs.StringVar(&c.host, "h", "localhost", "Host where thet server will listen for incoming requests")
	fs.Usage = func() {
		var usageString = `ccdns is a DNS Forwarder used to resolve DNS queries instead of directly using the authoritative nameserver chain.
		
		Usage of %s: <options> [value]`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return c, err
	}
	return c, nil
}

func runCmd(w io.Writer, c config) error {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.host, c.port))

	if err != nil {
		return err
	}

	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		<-signalCh
		// Initiate graceful shutdown
		fmt.Println("Bye bye ;-)")
		os.Exit(0)
	}()

	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	if err != nil {
		return err
	}

	// Read from UDP listener in endless loop
	for {
		var buf [512]byte
		_, _, err := conn.ReadFromUDP(buf[0:]) // adr
		if err != nil {
			return err
		}

		// MSG received
		// fmt.Println("> ", string(buf[0:]))
		fmt.Println("MSG received")
		// Write back the message over UPD
		//conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
	}
}
