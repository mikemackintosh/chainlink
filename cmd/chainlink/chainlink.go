package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/mikemackintosh/chainlink/config"
)

var (
	flagConfig            string
	flagListenDNS         int
	flagListenHTTP        int
	flagListenHTTPS       int
	flagChangeSysSettings string
)

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&flagConfig, "c", "config.yaml", "Configuration file")
	flag.StringVar(&flagChangeSysSettings, "i", "", "Changes system settings (Resolvers)")
	flag.IntVar(&flagListenDNS, "listen-dns", 53, "DNS Listen Port")
	flag.IntVar(&flagListenHTTP, "listen-http", 80, "HTTP Listen Port")
	flag.IntVar(&flagListenHTTPS, "listen-https", 443, "HTTPS Listen Port")
}

func main() {
	flag.Parse()

	// Create the waitgroup
	var wg = &sync.WaitGroup{}

	configLogger := NewLogger("CONFIG_1", "\033[38;5;214m")
	configLogger.Printf("-=-=-=-=-=-=-\n")
	configLogger.Printf("  Chainlink\n")

	configLogger.Printf("Loading configuration file: %s\n", flagConfig)
	configLogger.Printf("-=-=-=-=-=-=-\n")
	if err := config.ParseFromFile(flagConfig); err != nil {
		configLogger.Fatalf(err.Error())
	}

	// Start the DNS Server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		lg := NewLogger("   DNS_2", "\033[38;5;56m")
		lg.Printf("Starting DNS Server\n")
		lg.Printf("#=> Port: %d\n", flagListenDNS)
		srv := &dns.Server{Addr: ":" + strconv.Itoa(flagListenDNS), Net: "udp"}
		srv.Handler = &dnsHandler{
			logger:   lg,
			upstream: config.Registry.DNS.Upstream,
			zones:    config.Registry.GetResolvers(),
		}

		if err := srv.ListenAndServe(); err != nil {
			lg.Printf("Failed to set udp listener %s\n", err.Error())
		}
	}(wg)

	// Start the DNS Server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		lg := NewLogger("  HTTP_3", "\033[38;5;75m")
		lg.Printf("Starting HTTP Server\n")
		lg.Printf("#=> Port: %d\n", flagListenHTTP)

		srv := httpServer{
			logger: lg,
			server: &http.Server{
				Addr: fmt.Sprintf(":%d", flagListenHTTP),
			},
			Routes: config.Registry.GetRoutes(),
		}
		srv.server.Handler = srv.New()
		srv.server.RegisterOnShutdown(func() {
			wg.Done()
		})

		if err := srv.server.ListenAndServe(); err != nil {
			lg.Printf("Failed to set HTTP listener %s\n", err.Error())
		}
	}(wg)

	// Start the DNS Server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		lg := NewLogger(" HTTPS_4", "\033[38;5;61m")
		lg.Printf("Starting HTTPS Server\n")
		lg.Printf("#=> Port: %d\n", flagListenHTTPS)

		srv := httpServer{
			logger: lg,
			server: &http.Server{
				Addr: fmt.Sprintf(":%d", flagListenHTTPS),
			},
			Routes: config.Registry.GetRoutes(),
		}
		srv.server.Handler = srv.New()
		srv.server.RegisterOnShutdown(func() {
			wg.Done()
		})

		if err := srv.ListenAndServeTLS(); err != nil {
			lg.Printf("Failed to set HTTPS listener %s\n", err.Error())
		}
	}(wg)

	wg.Wait()
}
