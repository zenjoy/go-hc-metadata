A Go client to interact with the Hetzner Cloud Metadata API

This has been adapted from https://github.com/digitalocean/go-metadata/ to be used within the Hetzner Cloud.

# Usage

```go
// Create a client
client := metadata.NewClient(opts)

// Request all the metadata about the current droplet
all, err := client.Metadata()
if err != nil {
    log.Fatal(err)
}

// Lookup what our IPv4 address is on our first public
// network interface.
publicIPv4Addr := all.NetworkConfig.Config[0].Address

fmt.Println(publicIPv4Addr)
```

# License

MIT license
