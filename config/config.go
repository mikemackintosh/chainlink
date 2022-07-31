package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var Registry *Configuration

type Configuration struct {
	DNS   *UpstreamService `yaml:"dns"`
	Mgmt  *UpstreamService `yaml:"mgmt"`
	HTTPS ServerHTTPS      `yaml:"https"`
	Zones []*Zone          `yaml:"zones"`
}

type UpstreamService struct {
	Upstream string `yaml:"upstream"`
}

type ServerHTTPS struct {
	TLSCert string `yaml:"tls_cert"`
	TLSKey  string `yaml:"tls_key"`
}

type Resolve struct {
	Fqdn  string `yaml:"-"`
	Type  string `yaml:"type"` // currently ignored in the code
	Value string `yaml:"value"`
	TTL   uint32 `yaml:"ttl"`
}
type Route struct {
	Fqdn     string            `yaml:"-"`
	Path     string            `yaml:"path"`
	Upstream string            `yaml:"upstream"`
	Headers  map[string]string `yaml:"headers"`
}
type Endpoint struct {
	Resolve Resolve `yaml:"resolve"`
	Route   Route   `yaml:"http"`
}
type Zone struct {
	Zone      string               `yaml:"zone"`
	Endpoints map[string]*Endpoint `yaml:"endpoints"`
}

// ParseFromFile will read the provided config file path and decode it into a struct.
func ParseFromFile(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	c, err := Decode(b)
	if err != nil {
		return err
	}

	Registry = &c

	return nil
}

// SaveFile will save the provided config file path.
func SaveFile(f string) error {
	b, err := json.Marshal(Registry)
	if err != nil {
		return err
	}

	c, err := os.Create(f)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Decode will decode the configuration file bytes into the Configuration struct.
func Decode(b []byte) (Configuration, error) {
	var c Configuration
	err := yaml.Unmarshal([]byte(b), &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// GetResolvers will loop through all zones and endpoint and return all configurred
// resolvers as []*Resolve.
func (r *Configuration) GetResolvers() []*Resolve {
	var zones = []*Resolve{}

	// Add an internal route
	zones = append(zones, &Resolve{Fqdn: "chainlink.config.", Type: "A", Value: "127.0.0.1", TTL: 3600})

	for _, zone := range r.Zones {
		for name, endpoint := range zone.Endpoints {
			endpoint.Resolve.Fqdn = fmt.Sprintf("%s.%s.", name, zone.Zone)
			endpoint.Resolve.Type = strings.ToUpper(endpoint.Resolve.Type)
			zones = append(zones, &endpoint.Resolve)
		}
	}

	return zones
}

// GetRoutes will loop through all zones and endpoint and return all configurred
// routes as []*Route.
func (r *Configuration) GetRoutes() []*Route {
	var zones = []*Route{}

	zones = append(zones, &Route{Fqdn: "chainlink.config.", Path: "/*", Upstream: "http://" + r.Mgmt.Upstream, Headers: map[string]string{}})

	for _, zone := range r.Zones {
		for name, endpoint := range zone.Endpoints {
			endpoint.Route.Fqdn = fmt.Sprintf("%s.%s.", name, zone.Zone)
			zones = append(zones, &endpoint.Route)
		}
	}

	for _, z := range zones {
		fmt.Printf("%#v\n", *z)
	}

	return zones
}
