package main

import (
	"fmt"
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
			switch zone.Type {
			case "A":
				msg.Answer = append(msg.Answer, &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    zone.TTL,
					},
					A: net.ParseIP(zone.Value),
				})
			case "AAAA":
				msg.Answer = append(msg.Answer, &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    zone.TTL,
					},
					AAAA: net.ParseIP(zone.Value),
				})
			case "CNAME":
				msg.Answer = append(msg.Answer, &dns.CNAME{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeCNAME,
						Class:  dns.ClassINET,
						Ttl:    zone.TTL,
					},
					Target: zone.Value,
				})
			}
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
	h.logger.Printf("Querying %v %s...\n", dns.Type(queryType).String(), queryDomain)
	h.ResolveQuery(queryType, queryDomain, &msg)

	// Implies 0 = no response, so lets proxy upstream.
	if len(msg.Answer) == 0 {
		h.logger.Printf("-> Not found locally, sending upstream\n")
		if err := h.proxyUpstream(w, req); err != nil {
			h.logger.Printf("Err with upstream: %s\n", err)
			dns.HandleFailed(w, req)
		}

		return
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
	req.RecursionDesired = true
	resp, _, err := c.Exchange(req, h.upstream)
	if err != nil {
		return err
	}

	for _, answer := range resp.Answer {
		fmt.Printf("\033[38;5;214m[ANSWER]\033[0m %#v\n", answer)
	}

	if err := w.WriteMsg(resp); err != nil {
		return err
	}

	return nil
}
