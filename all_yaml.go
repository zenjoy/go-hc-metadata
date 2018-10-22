package metadata

type Metadata struct {
	PublicIPv4 string   `yaml:"public-ipv4,omitempty"`
	LocalIPv4  string   `yaml:"local-ipv4,omitempty"`
	Hostname   string   `yaml:"hostname,omitempty"`
	PublicKeys []string `yaml:"public-keys,omitempty"`
	InstanceId int   `yaml:"instance-id,omitempty"`
	VendorData string   `yaml:"vendor_data,omitempty"`
	NetworkSysconfig string   `yaml:"network-sysconfig,omitempty"`

	NetworkConfig struct {
		Version string   `yaml:"version,omitempty"`
		Config []struct {
			MACAddress string `yaml:"mac_address,omitempty"`
			Type       string `yaml:"type,omitempty"`
			Name       string `yaml:"name,omitempty"`
			
			Subnets []struct {
				DNSNameServers []string `yaml:"dns_nameservers,omitempty"`
				Type  string `yaml:"name,omitempty"`
				IPv4  bool `yaml:"ipv4,omitempty"`
				IPv6  bool `yaml:"ipv6,omitempty"`
				Address  string `yaml:"address,omitempty"`
				Gateway  string `yaml:"gateway,omitempty"`

				Routes []struct {
					Netmask  int `yaml:"netmask,omitempty"`
					Gateway  string `yaml:"gateway,omitempty"`
					Network  string `yaml:"network,omitempty"`
				} `yaml:"routes,omitempty"`
			} `yaml:"subnets,omitempty"`
		} `yaml:"config,omitempty"`
	} `yaml:"network-config",omitempty"`
}
