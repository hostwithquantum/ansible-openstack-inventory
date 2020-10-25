package auth

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"

	"github.com/gophercloud/gophercloud/openstack"
	"github.com/pkg/errors"
)

// Authenticate against Keystone and create a client
func Authenticate() (*gophercloud.ProviderClient, error) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pc, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tlsconfig := &tls.Config{}
	tlsconfig.InsecureSkipVerify = true
	transport := &http.Transport{TLSClientConfig: tlsconfig}
	pc.HTTPClient = http.Client{
		Transport: transport,
	}

	err = openstack.Authenticate(pc, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pc, nil
}
