package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/miekg/dns"
	"github.com/mikemackintosh/chainlink/config"
)

var (
	flagConfig            string
	flagChangeSysSettings bool
	flagListenDNS         int
	flagListenHTTP        int
	flagListenHTTPS       int

	httplog = NewLogger("  HTTP_3", "\033[38;5;75m")
)

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&flagConfig, "c", "config.yaml", "Configuration file")
	flag.BoolVar(&flagChangeSysSettings, "i", false, "Changes system settings (Resolvers)")
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

		httplog.Printf("Starting HTTP Server\n")
		httplog.Printf("#=> Port: %d\n", flagListenHTTP)

		srv := httpServer{
			logger: httplog,
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
			httplog.Printf("Failed to set HTTP listener %s\n", err.Error())
		}
	}(wg)

	// Start the HTTPS Server
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

	// Start the Management Server
	if config.Registry.Mgmt != nil {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			lg := NewLogger(" MGMT_5", "\033[38;5;86m")
			lg.Printf("Starting Management Server\n")
			lg.Printf("#=> Port: %s\n", config.Registry.Mgmt.Upstream)

			// Set the handler to output the request
			r := chi.NewRouter()
			// Basic CORS
			// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
			r.Use(cors.Handler(cors.Options{
				// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
				AllowedOrigins: []string{"https://chainlink.dev", "http://localhost:3000", "http://127.0.0.1:3000"},
				// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300, // Maximum value not ignored by any of major browsers
			}))
			r.Use(middleware.RequestID)
			r.Use(middleware.RealIP)
			r.Use(middleware.Logger)
			r.Use(middleware.Recoverer)

			// Set the default route for React
			r.Get("/*", handlerReactProxy)

			// Configure the api route mount
			r.Mount("/api/", apiRouter())

			log.Fatal(http.ListenAndServe(config.Registry.Mgmt.Upstream, r))

		}(wg)
	}

	wg.Wait()
}

// A completely separate router for administrator routes
func apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/config", adminIndex)
	r.Post("/config", adminUpdate)
	r.Delete("/config", adminDelete)
	return r
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	b, _ := json.MarshalIndent(config.Registry.Zones, "  ", "  ")

	w.Write(b)
}

type ZoneUpdateRequest struct {
	Fqdn     string `json:"fqdn"`
	Upstream string `json:"upstream"`
}

func adminUpdate(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var newZone ZoneUpdateRequest
	err = json.Unmarshal(b, &newZone)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	u, err := url.Parse(newZone.Fqdn)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	parts := strings.Split(u.Path, ".")
	zone := parts[len(parts)-2] + "." + parts[len(parts)-1]
	hostname := strings.Replace(u.Path, "."+zone, "", -1)

	endpoint := &config.Endpoint{
		Resolve: config.Resolve{Fqdn: u.Path + ".", Type: "A", Value: "127.0.0.1", TTL: 3600},
		Route:   config.Route{Fqdn: u.Path + ".", Path: "/*", Upstream: newZone.Upstream, Headers: map[string]string{}},
	}

	ep := map[string]*config.Endpoint{hostname: endpoint}
	if _, z, ok := config.Registry.FindZone(zone); ok {
		z.Endpoints[hostname] = endpoint
	} else {
		config.Registry.Zones = append(config.Registry.Zones, &config.Zone{Zone: zone, Endpoints: ep})
	}

	var host = strings.Trim(endpoint.Route.Fqdn, ".")
	httplog.Printf("Configuring %s -> %s\n", host+endpoint.Route.Path, endpoint.Route.Upstream)

	router := chi.NewRouter()
	p := Proxy{Upstream: endpoint.Route.Upstream, Path: endpoint.Route.Path, Headers: endpoint.Route.Headers, logger: httplog}
	router.HandleFunc(p.Path, p.ProxyHandler)

	hr.Map(host, router)

	j, err := json.Marshal(map[string]interface{}{"zone": zone, "endpoints": ep})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(j)
}

func adminDelete(w http.ResponseWriter, r *http.Request) {
	b, _ := json.MarshalIndent(config.Registry.Zones, "  ", "  ")

	w.Write(b)
}
