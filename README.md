# chainlink
---

Chainlink is a development resource designed to isolate your local services via DNS and HTTP routing.

**Use case**: You are developing locally for multiple services and you need to receive requests without modifying your `/etc/resolv.conf`.

I use Chainlink for my all of local development. This tool allows me to impersonate production services, which I can run locally, and reverse proxy the requests I need to specific HTTP and DNS upstreams. In the event a DNS entry is not available, it will resolve by default against `8.8.8.8:53`.

## Installation

```bash
# Download the package to src/
go get github.com/mikemackintosh/chainlink

# Install the package
go install github.com/mikemackintosh/chainlink

# Copy the stub config file
cp $GOPATH/src/github.com/mikemackintosh/chainlink/{example-config.yaml,config.yaml}

# Edit the configuration file
$EDITOR $GOPATH/src/github.com/mikemackintosh/chainlink/config.yaml
```

## Usage

Once installed (below), execute chainlink using the following command:

```bash
$ chainlink
```

It expects a local configuration file at `config.yaml`:
```yaml
---
dns:
  upstream: 8.8.8.8:53 # Change to your local DNS if desired.
zones:
  - zone: "zush.int" # Tld for resolution
    endpoints:

      # accounts.zush.int => localhost:8081
      accounts:
        resolve:
          type: a
          value: 127.0.0.1
          ttl: 600
        http:
          path: /*
          upstream: http://localhost:8081
          headers:
            Authorization: "Bearer blahblah"
```
