package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/miekg/dns"
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

func parseArgs(w io.Writer, args []string) (*config, error) {
	c := config{}
	fs := flag.NewFlagSet("ccdns", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.IntVar(&c.port, "p", 53, "Port where thet server will listen for incoming requests")
	fs.StringVar(&c.host, "h", "", "Host where thet server will listen for incoming requests")
	fs.Usage = func() {
		var usageString = `ccdns is a DNS Forwarder used to resolve DNS queries instead of directly using the authoritative nameserver chain.
		
		Usage of %s: <options> [value]`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return &c, err
	}
	return &c, nil
}

func runCmd(w io.Writer, c *config) error {
	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		<-signalCh
		// Initiate graceful shutdown
		fmt.Println("Bye bye ;-)")
		os.Exit(0)
	}()

	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      fmt.Sprintf("%s:%d", c.host, c.port),
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Printf("Starting DNS server on %v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
	return nil

}

type dnsHandler struct{}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		//fmt.Printf("Received query: %s\n", question.Name)
		answers := resolve(question.Name, question.Qtype)
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
}

func resolve(domain string, qtype uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return in.Answer
}
