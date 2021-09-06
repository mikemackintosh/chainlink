package main

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/mikemackintosh/chainlink/config"
)

const DefaultUpstreamDNS = "8.8.8.8:53"

type dnsHandler struct {
	logger   *Logger
	upstream string
	zones    []*config.Resolve
}

func (h *dnsHandler) ResolveQuery(queryType uint16, domain string, msg *dns.Msg) {
	d := strings.ToLower(domain)
	for _, zone := range h.zones {
		if zone.Fqdn == d {
			// TODO: record type switching for a later version
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: zone.TTL},
				A:   net.ParseIP(zone.Value),
			})
		}
	}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	var msg = dns.Msg{}
	msg.SetReply(req)

	// Get the query type and domain name from the request
	if len(req.Question) < 1 {
		h.logger.Println("Err with request, not a well structured query")
		dns.HandleFailed(w, req)
		return
	}

	var queryType = req.Question[0].Qtype
	var queryDomain = msg.Question[0].Name
	h.logger.Printf("Querying %v %s...\n", queryType, queryDomain)
	h.ResolveQuery(queryType, queryDomain, &msg)

	// Implies 0 = no response, so lets proxy upstream.
	if len(msg.Answer) == 0 {
		h.logger.Printf("-> Not found locally, sending upstream\n")
		if err := h.proxyUpstream(w, req); err != nil {
			h.logger.Printf("Err with upstream: %s\n", err)
			dns.HandleFailed(w, req)
			return
		}
	}

	if err := w.WriteMsg(&msg); err != nil {
		h.logger.Printf("Err: %s\n", err)
	}
}

// proxyUpstream will proxy all requests to the configured server.
func (h *dnsHandler) proxyUpstream(w dns.ResponseWriter, req *dns.Msg) error {
	var transport = "udp"
	if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		transport = "tcp"
	}

	if len(h.upstream) == 0 {
		h.upstream = DefaultUpstreamDNS
	}

	var c = &dns.Client{Net: transport}
	resp, _, err := c.Exchange(req, h.upstream)
	if err != nil {
		return err
	}
	if err := w.WriteMsg(resp); err != nil {
		return err
	}

	return nil
}
