package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/mikemackintosh/chainlink/config"
)

type httpServer struct {
	server *http.Server
	logger *Logger
	Routes []*config.Route
}

// New will create a new httpServer and return the chi.Mux equivelant. This is to then be passed
// to a http.Server{Handler: xxx }.
func (s httpServer) New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Update the default chi logging middleware to match our output
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: s.logger})
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create a host router
	hr := hostrouter.New()
	for _, endpoint := range s.Routes {
		var host = strings.Trim(endpoint.Fqdn, ".")
		s.logger.Printf("Configuring %s -> %s\n", host+endpoint.Path, endpoint.Upstream)

		r := chi.NewRouter()
		p := Proxy{Upstream: endpoint.Upstream, Path: endpoint.Path, Headers: endpoint.Headers, logger: s.logger}
		r.HandleFunc(p.Path, p.ProxyHandler)

		hr.Map(host, r)
	}

	// Mount the host router
	r.Mount("/", hr)

	return r
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// ListenAndServeTLSKeyPair start a server using in-memory TLS KeyPair
func (s httpServer) ListenAndServeTLS() error {
	var subjectAltNames = []string{}
	for _, endpoint := range s.Routes {
		subjectAltNames = append(subjectAltNames, strings.Trim(endpoint.Fqdn, "."))
	}

	cert, err := GenTLSKeyPair(subjectAltNames)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)},
		&tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: []tls.Certificate{cert},
		})

	return s.server.Serve(tlsListener)
}

// Proxy used to pass details from the configurator to the Proxy handler.
type Proxy struct {
	logger   *Logger
	Upstream string
	Path     string
	Headers  map[string]string
}

// ProxyHandler is an HTTP proxy handler for specific routes. It will upstream the request
// and update the ResponseWriter accordingly.
func (p Proxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(p.Upstream)
	proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
		r.URL.Scheme = url.Scheme
		r.URL.Host = url.Host
		r.URL.Path = url.Path + r.URL.Path
		r.Host = url.Host
		for k, v := range p.Headers {
			r.Header.Set(k, v)
		}
	}}
	proxy.ServeHTTP(w, r)
}

// GenTLSKeyPair generates the TLS keypair for the server.
func GenTLSKeyPair(hostnames []string) (tls.Certificate, error) {
	now := time.Now()

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "localhost",
			Country:      []string{"USA"},
			Organization: []string{"mikemackintosh"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0), // Valid for one day
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              hostnames,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	//	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template,
		priv.Public(), priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	var outCert tls.Certificate
	outCert.Certificate = append(outCert.Certificate, cert)
	outCert.PrivateKey = priv

	return outCert, nil
}
