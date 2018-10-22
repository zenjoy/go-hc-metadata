package metadata

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
  "testing"
  "io/ioutil"
  
  "gopkg.in/yaml.v2"
)

func TestInstanceID(t *testing.T) {
	var (
		resp = "4567"
		want = 4567
	)
	withServer(t, "/2009-04-04/meta-data/instance-id", resp, func(client *Client) {
		got, err := client.InstanceID()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf(" got=%#v", got)
		}
	})
}

func TestHostname(t *testing.T) {
	var (
		resp = "localhost"
		want = "localhost"
	)
	withServer(t, "/2009-04-04/meta-data/hostname", resp, func(client *Client) {
		got, err := client.Hostname()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf(" got=%#v", got)
		}
	})
}

func TestUserData(t *testing.T) {
	var (
		resp = "#!/bin/sh\necho 'hello world'"
		want = "#!/bin/sh\necho 'hello world'"
	)
	withServer(t, "/hetzner/v1/userdata", resp, func(client *Client) {
		got, err := client.UserData()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf(" got=%#v", got)
		}
	})
}

func TestVendorData(t *testing.T) {
	var (
		resp = "#!/bin/sh\necho 'hello world'"
		want = "#!/bin/sh\necho 'hello world'"
	)
	withServer(t, "/2009-04-04/meta-data/vendor_data", resp, func(client *Client) {
		got, err := client.VendorData()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf(" got=%#v", got)
		}
	})
}

func TestPublicKeys(t *testing.T) {
	var (
		resp = "ssh-rsa sshkeysshkeysshkey1 user@workstation2\nssh-rsa sshkeysshkeysshkey2 user@workstation2"
		want = []string{"ssh-rsa sshkeysshkeysshkey1 user@workstation2", "ssh-rsa sshkeysshkeysshkey2 user@workstation2"}
	)
	withServer(t, "/2009-04-04/meta-data/public-keys", resp, func(client *Client) {
		got, err := client.PublicKeys()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf(" got=%#v", got)
		}
	})
}

func TestMetadata(t *testing.T) {
	resp := `hostname: hc-node-0
instance-id: 1234567
local-ipv4: ''
network-config:
  config:
  - mac_address: DE:AD:BE:EF:DE:AD
    name: eth0
    subnets:
    - dns_nameservers:
      - 213.133.100.100
      - 213.133.98.98
      - 213.133.99.99
      ipv4: true
      type: dhcp
    type: physical
  - name: eth0:0
    subnets:
    - address: 2a01:4f8:c010:14ab::1/64
      gateway: fe80::1
      ipv6: true
      routes:
      - gateway: fe80::1%eth0
        netmask: 0
        network: '::'
      type: static
    type: physical
  version: 1
network-sysconfig: "DEVICE='eth0'\nTYPE=Ethernet\nBOOTPROTO=dhcp\nONBOOT='yes'\nHWADDR=DE:AD:BE:EF:DE:AD\n\
  IPV6INIT=yes\nIPV6ADDR=2a01:4f8:c010:14ab::1/64\nIPV6_DEFAULTGW=fe80::1%eth0\nIPV6_AUTOCONF=no\n\
  DNS1=213.133.100.100\nDNS2=213.133.98.98\n        "
public-ipv4: 159.69.180.192
public-keys:
- 'ssh-rsa sshkeysshkeysshkey1 user@workstation'
vendor_data: "Content-Type: multipart/mixed; boundary=\"===============5184703028565742311==\"\
  \nMIME-Version: 1.0\n\n--===============5184703028565742311==\nContent-Type: text/text/cloud-config;\
  \ charset=\"us-ascii\"\nMIME-Version: 1.0\nContent-Transfer-Encoding: 7bit\nContent-Disposition:\
  \ attachment; filename=\"cloud-config\"\n\n#cloud-config\nfqdn: flynn-node-0\nmanage_etc_hosts:\
  \ true\nrandom_seed:\n  data: !!binary |\n    L2gxUE1wSlkzNVczUm94aDZFdzVlcmt5YTRVSW8xNFNIQnpMaUczWStMZTI3Z3FnWEhBclVYTzNM\n\
  \    Q1Zjd0ZUemErcytFQUt0T29CRDFRTGR6dFRYUktXVnFVK0RzKzY4QlR0Q0R5cUErRVg2U0hka1g5\n\
  \    K09QYU1BVkYxbytRVXFOTWxtOW5lanc2VS9tOGJlV3BhYVd1QXF0ZDR3RG5jYXliN1FCaE02RmVU\n\
  \    SmVTeDNqRTdBejNzejYrdjh0c2RENUZpeDl4bzZmVEhBaEQxTzhlV0RQN2p4U1kxTHMrK3VmbTNZ\n\
  \    cXBQcnN4eG1JMjNDcUhndkcyaDdrQitqNGduUEhFbXhrQzJFMmRGbFBYY3hra3hWMVRQd3FoL0Y0\n\
  \    QkZURk1oYjVMdGpiNE54RmU3TEpKQ1BuZ3FUNG1WM09XbmIxL1M0QU81N1ZuV1NiUkFEdkpaUXhr\n\
  \    Vjg4OUdWa1l3Q0VTaDdXSTFUR1JVNmhuckhaNy8vVHN4WHJkbW44UUNQK1ZSOExlTDZRbzF5Zmgz\n\
  \    cm92bXVlVVVMRzhhMm5xc01UQUcyRFczRzZYRFYxS1I1ajN5dmtUcFZ3Tk5qaytoV0FZWEczTXYx\n\
  \    dzFCd1Rid2xKTFQ3ZllSWXYxSjRERE9rT0c1ZHJPVXNJRzU4bzZhVHRsUXkzc3RBeE5MWTM4d2ta\n\
  \    U1p6VVFLTFlVRmFwbjlYeVcyNTY4L3V2ZzI3RXZNTHB2cWxEK1g0QkRBWVVPQzQrUXV2bGlJYWUx\n\
  \    Q0Z4OG1IUUMySjBuNFFERnY5bmJqQTgyaGFLcG9teFNGVjVSWGRpeGxYN1RMN0NzTldQbzZIMDNR\n\
  \    VlhCQ1JUMXp2bnVjMlRrOERXbS9YUmZGUXVOSmd2aHdURXRoTmJxcmVCaXZDRGcxb0ZUTmxwSlVu\n\
  \    NGZLUVNEK252bVhwTElseVJEUFZSYkJwQmJhbFB6TjRkTUx6OTY2MjcrRXNRL3NUcmZ1RXRIaEto\n\
  \    clBxWmFUN25YOXdRTFpsUGppOFlGeU5KN2t2S1VKKzM4dGV4WmIyUytTVDRkQUJUZ1d5NEZhdExS\n\
  \    eXhpMW14bjY5UkM3UjlrSy91TG1SSW15VWxQQnZPZ2hpZHorRzBRYXFxdEloSHhub1FzRkYvdkdB\n\
  \    Rlg2UlQyZlVsanlYUEJDYkh3RVBVRkxNbTQza0hOUU90Mi9lbmpqWFlobFB3TWppUWZEMVRPUzBE\n\
  \    UHBvcHMrVWtDbGp3alFLZERkaDNieDI0UWs2OEtjV2JwOUpCamJKVlJ5Sk4xeUhMeDZieWV5OVVH\n\
  \    WVVud0NrTE5BbnJqTkxONXVSK0hCMTQrRjcrVnNtU3l2dHdPYjN3Mk1Ma1Q0VXNCaUdlb21SMXMy\n\
  \    aWw4SjhIZXlHY2NLOUhPTU45Unh6N004OWpOa3FacTBxK1IyZ1JNQmJ2ZjFEVU96YmZrNGg1MktY\n\
  \    dU1veC9MNWR4RG9rRXplUEhEQjV1dXMrWW1MQTRKblJXUHFDd0ZNY2NxTmY0aVZFQzhaYlNLZS83\n\
  \    Rm1oSnRhYVQyYTFPOWtlR1cvNlVHZW9HTE44UVR2VzBZWUhGajRuZ1Irbm5rWGNXblBGUXpyMWRR\n\
  \    RHBxTjNZbDR6NWprZkErVHcydEpBMDE1WmdZUEhGMFcrY3VON1MwdityWEd0Wm5ML0hyZWxzK3dh\n\
  \    M2RNelVnNWRRQ1FDNnZTYjA4TG00NTYrWnBJYWtwOWkvMmQwMDdOb01Ra1VFWlZLT05lSjA2bFRu\n\
  \    cWFnNHdVYVdGZFdkNllIM29GRFhaWWlyUTZXNXArVmNkSVVMbW50ZHRyTk52dFd0NGJpQzN4dGVH\n\
  \    TWdGc1dYYUZJR0FLSWNPb0pLNWNsZFNMZlNZbWhydVFLSmk1NUVHdmJvL2Q1ZEIzaTM1Tm9ocU9J\n\
  \    REczYVdwb1FTcWV1TDByMzdRelg4WDRVZ1hyVGJPZy93bUR1NTZSaDFJYmZCUEplODdOV3hqODBs\n\
  \    NmVmOUtHQ3dqNzByRndPMmFRTm9vRzlhZmxlc2dpSGJZWU5qNWhHZlUxd214UTFRcjFvMkJLVk5q\n\
  \    c3ZTckF1ek04UWZraFFxYjM2TktoRXhYc2E3SWs5aFdrM2l5T2JjM3NESzZkZXZGVndjVEZxQ2FO\n\
  \    cUdUWWZUeSszMW4za201Z3NWOGRVSHlXREFUNng5eTRiQnB6VjFkSks1b0Z1NTVLYm10QXczSFBB\n\
  \    akZ3SGQ0YnVGMFVrUlBKYU53SU1lU0ZJN3lwcS8rYVdCNzUzVnVZbXQ2a2UxNE9CNHBXWG1QbFlR\n\
  \    Rkh4MGFXTUhjbzBmOTdTbUlEeG9QVVZOQ1FXOU1CSDdUVjdWU09uUUlTQ0twVGtqcXk3Y1VIaTJH\n\
  \    WjlpL3FGUlhpSDZSSXlsM20xUkFwSEUvRytYVGxMRFdFOGY1cWwwcURBbzE3YzJIOHZSR0ZVVCtz\n\
  \    a3VoblJCcXUwQTkrdDNBak1NRVE2RjlhRUJWR0hiTGZNQzhqdWlkU05qQW1UazdFMTUwbHdZK3pn\n\
  \    NGFNUjBDdmIvMXk0YWZaT3pqQWtWd1ovY2huUzdpbnVlbXU4YUROZUp1c1JHSm1WRTEveW5LTXRC\n\
  \    a2NHc2tDWE5VS1MySTJmK2RUQVQ3VnpoQzg0UERTR1VhNm5NSkRibWlhaUtHS3BDQ3o1b0E2M2M2\n\
  \    UUJFd1UvVFBDbDRIaThvcmdjTnpIeFJ4NFZ4eGdZNDYzN25xWnMva0ZUS2w5TUpTR1FEL0FvTjVK\n\
  \    TW1lV1Z5QjZoRE9IbjRULzNoMzBBdkg0ejlwV0lkd0xjSTQ4bTFWQ3VPNkxuWWZvVFhwNm9tZW10\n\
  \    N2lsNDdWQkxQR0FNaUdpQ3oxYThjVGpiek9sSlNLS0JyRjlmSXZBbW9GTzdJc1VISjZpMksyMlZO\n\
  \    NUJrNmhUbnlNYjd3dGVXb1dqazl5L2l4QUoyMEFKWkRkcUNKdFgrZm03RktELzRVTWlUT1lPWmt5\n\
  \    cU5ON2tzcFRJcFNrY3lCMHhTSDBsd25jWFdBcWlpcS80aGdob2xuek9OY1huTjJXcGY4bjlTa1hC\n\
  \    V3A0N1RhYzZYNVRFVW1PemhvekRaUllyNTZBbExaTWFCb1VxYnRReXY3RDRCUHRNWnNKQjMrSkt5\n\
  \    RjdQZ0xMQTdYZmhpMHpZYUlTOTc0QTV4VGU2SUtyTjM5WmZXVW9HYi94ZXE1V3UyZ3pYZEZXSW10\n\
  \    bE5LRUNpNVNvRk5VcDZobWZ5MTBwSXFXSHU3T0ZIOWE5OEVFTmdFYXRWUjErY0g2QzBTTW55ZFlE\n\
  \    Z2xUaExtdC9IODdaN1NFRmZ4RkVwa2luNEgvTnZxZ3pMY1gyRmhKNkQ3c1JIQ3B0T1ZQMi9ncWcw\n\
  \    VHl5OXF0Sytpemt1SnlISjgzcGFxQ1h4dkpMWXZFZXJPRFZmR0FFL2RFZ053dkNrQ3BVclZteURi\n\
  \    VDMvczNzS3NNM2tIdW1ZTGhLSDhnUEdzcUNPWm03RDJSQzNsMTh1RW1NQnhNSVhjQWFIUGloNENm\n\
  \    U1h1T0JOUlFCcThjMFBHMXp0VHJacFdDQzFFU1d1SjgxS3IrZnNTdmhiR1hDTUdraCswaWZlbE9E\n\
  \    MUtPQVRLRExZMGdPOEhWZHREYi9VM2tLYVRkVHRCQ01TS3NjMXV0ejlCb1E3d0ZmbnJxZz0=\n\
  \  encoding: base64\n  file: /dev/urandom\nsystem_info:\n  default_user:\n    lock_passwd:\
  \ true\n    name: root\n    shell: /bin/bash\n  distro: ubuntu\n\n--===============5184703028565742311==\n\
  Content-Type: text/text/cloud-boothook; charset=\"us-ascii\"\nMIME-Version: 1.0\n\
  Content-Transfer-Encoding: 7bit\nContent-Disposition: attachment; filename=\"hc-boot-script\"\
  \n\n#!/bin/bash\n\n--===============5184703028565742311==--\n"`

	var (
		want = new(Metadata)
		got  *Metadata
	)

  var respBytes []byte

  respBytes, _ = ioutil.ReadAll(strings.NewReader(resp))

	if err := yaml.Unmarshal(respBytes, &want); err != nil {
		t.Fatalf("%#v", err)
	}

	withServer(t, "/hetzner/v1/metadata", resp, func(client *Client) {
		var err error
		got, err = client.Metadata()
		if err != nil {
			t.Fatal(err)
		}
	})
	if !reflect.DeepEqual(*want, *got) {
		t.Errorf("want=%#v", *want)
		t.Errorf(" got=%#v", *got)
	}
}

func withServer(t testing.TB, path, response string, test func(*Client)) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.Error(rw, "bad path", http.StatusBadRequest)
			t.Errorf("bad URL sent to server: %v", req.URL.String())
			return
		}
		rw.Write([]byte(response))
	}))
	defer srv.Close()
	u, err := url.Parse(srv.URL)
	if err != nil {
		panic(err)
	}
	test(NewClient(WithBaseURL(u)))
}
