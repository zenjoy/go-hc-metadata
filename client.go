// Package metadata implements a client for the DigitalOcean metadata
// API. This API allows a droplet to inspect information about itself,
// like it's region, droplet ID, and so on.
//
// Documentation for the API is available at:
//
//    https://developers.digitalocean.com/documentation/metadata/
package metadata

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	maxErrMsgLen = 128 // arbitrary max length for error messages

	defaultTimeout = 2 * time.Second
	defaultPath    = "/2009-04-04/meta-data"
)

var (
	defaultBaseURL = func() *url.URL {
		u, err := url.Parse("http://169.254.169.254")
		if err != nil {
			panic(err)
		}
		return u
	}()
)

// ClientOption modifies the default behavior of a metadata client. This
// is usually not needed.
type ClientOption func(*Client)

// WithHTTPClient makes the metadata client use the given HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(metaclient *Client) { metaclient.client = client }
}

// WithBaseURL makes the metadata client reach the metadata API using the
// given base URL.
func WithBaseURL(base *url.URL) ClientOption {
	return func(metaclient *Client) { metaclient.baseURL = base }
}

// Client to interact with the DigitalOcean metadata API, from inside
// a droplet.
type Client struct {
	client  *http.Client
	baseURL *url.URL
}

// NewClient creates a client for the metadata API.
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		client:  &http.Client{Timeout: defaultTimeout},
		baseURL: defaultBaseURL,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Metadata contains the entire contents of a Droplet's metadata.
// This method is unique because it returns all of the
// metadata at once, instead of individual metadata items.
func (c *Client) Metadata() (*Metadata, error) {
	metadata := new(Metadata)
	err := c.doGetURL(c.resolve("/hetzner/v1/metadata"), func(r io.Reader) error {
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(r)

		return yaml.Unmarshal(bodyBytes, &metadata)
	})
	return metadata, err
}

// InstanceID returns the Droplet's unique identifier. This is
// automatically generated upon Droplet creation.
func (c *Client) InstanceID() (int, error) {
	InstanceID := new(int)
	err := c.doGet("instance-id", func(r io.Reader) error {
		_, err := fmt.Fscanf(r, "%d", InstanceID)
		return err
	})
	return *InstanceID, err
}

// Hostname returns the Droplet's hostname, as specified by the
// user during Droplet creation.
func (c *Client) Hostname() (string, error) {
	var hostname string
	err := c.doGet("hostname", func(r io.Reader) error {
		hostnameraw, err := ioutil.ReadAll(r)
		hostname = string(hostnameraw)
		return err
	})
	return hostname, err
}

// UserData returns the user data that was provided by the user
// during Droplet creation. User data can contain arbitrary data
// for miscellaneous use or, with certain Linux distributions,
// an arbitrary shell script or cloud-config file that will be
// consumed by a variation of cloud-init upon boot. At this time,
// cloud-config support is included with CoreOS, Ubuntu 14.04, and
// CentOS 7 images on DigitalOcean.
func (c *Client) UserData() (string, error) {
	var userdata string
	err := c.doGetURL(c.resolve("/hetzner/v1/userdata"), func(r io.Reader) error {
		userdataraw, err := ioutil.ReadAll(r)
		userdata = string(userdataraw)
		return err
	})
	return userdata, err
}

// VendorData provided data that can be used to configure Droplets
// upon their creation. This is similar to user data, but it is
// provided by DigitalOcean instead of the user.
func (c *Client) VendorData() (string, error) {
	var vendordata string
	err := c.doGet("vendor_data", func(r io.Reader) error {
		vendordataraw, err := ioutil.ReadAll(r)
		vendordata = string(vendordataraw)
		return err
	})
	return vendordata, err
}

// PublicKeys returns the public SSH key(s) that were added to
// the Droplet's root user's authorized_keys file during Droplet
// creation.
func (c *Client) PublicKeys() ([]string, error) {
	var keys []string
	err := c.doGet("public-keys", func(r io.Reader) error {
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			keys = append(keys, scan.Text())
		}
		return scan.Err()
	})
	return keys, err
}

// Tags returns a list of DigitalOcean tags that have been
// applied to the Droplet.
func (c *Client) Tags() ([]string, error) {
	var ns []string
	err := c.doGet("tags", func(r io.Reader) error {
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			ns = append(ns, scan.Text())
		}
		return scan.Err()
	})
	return ns, err
}

func (c *Client) doGet(resource string, decoder func(r io.Reader) error) error {
	return c.doGetURL(c.resolve(defaultPath, resource), decoder)
}

func (c *Client) doGetURL(url string, decoder func(r io.Reader) error) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.makeError(resp)
	}
	return decoder(resp.Body)
}

func (c *Client) makeError(resp *http.Response) error {
	body, _ := ioutil.ReadAll(io.LimitReader(resp.Body, maxErrMsgLen))
	if len(body) >= maxErrMsgLen {
		body = append(body[:maxErrMsgLen], []byte("... (elided)")...)
	} else if len(body) == 0 {
		body = []byte(resp.Status)
	}
	return fmt.Errorf("unexpected response from metadata API, status %d: %s",
		resp.StatusCode, string(body))
}

func (c *Client) resolve(basePath string, resource ...string) string {
	dupe := *c.baseURL
	dupe.Path = path.Join(append([]string{basePath}, resource...)...)
	return dupe.String()
}
