/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package credentials implements various credentials supported by gRPC library,
// which encapsulate all the state needed by a client to authenticate with a
// server and make various assertions, e.g., about the client's identity, role,
// or whether it is authorized to make a particular call.
package credentials // import "google.golang.org/grpc/credentials"

import (
	//"crypto/tls"
	"github.com/gotlcp/tlcp"
	"fmt"
	//"crypto/x509"
	x509 "github.com/emmansun/gmsm/smx509"
	"io/ioutil"
	"net"
	"strings"

	"golang.org/x/net/context"
)

// TLCPInfo contains the auth information for a TLCP authenticated connection.
// It implements the AuthInfo interface.
type TLCPInfo struct {
	State tlcp.ConnectionState
}

// AuthType returns the type of TLCPInfo as a string.
func (t TLCPInfo) AuthType() string {
	return "tlcp"
}

// tlcpCreds is the credentials required for authenticating a connection using TLCP.
type tlcpCreds struct {
	// TLCP configuration
	config *tlcp.Config
}

func (c tlcpCreds) Info() ProtocolInfo {
	return ProtocolInfo{
		SecurityProtocol: "tlcp",
		SecurityVersion:  "1.1", // todo: liwei   1.1
		ServerName:       c.config.ServerName,
	}
}

func (c *tlcpCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) {
	// use local cfg to avoid clobbering ServerName if using multiple endpoints
	cfg := cloneTLCPConfig(c.config)
	if cfg.ServerName == "" { // todo liwei
		colonPos := strings.LastIndex(authority, ":")
		if colonPos == -1 {
			colonPos = len(authority)
		}
		cfg.ServerName = authority[:colonPos]
	}
	conn := tlcp.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- conn.Handshake()
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			return nil, nil, err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	return conn, TLCPInfo{conn.ConnectionState()}, nil
}

func (c *tlcpCreds) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) {
	conn := tlcp.Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	return conn, TLCPInfo{conn.ConnectionState()}, nil
}

func (c *tlcpCreds) Clone() TransportCredentials {
	return NewTLCP(c.config)
}

func (c *tlcpCreds) OverrideServerName(serverNameOverride string) error {
	c.config.ServerName = serverNameOverride
	return nil
}

// NewTLCP uses c to construct a TransportCredentials based on TLCP.
func NewTLCP(c *tlcp.Config) TransportCredentials {
	tc := &tlcpCreds{cloneTLCPConfig(c)}
	//tc.config.NextProtos = alpnProtoStr // todo liwei
	return tc
}

// NewClientTLCPFromCert constructs TLCP credentials from the input certificate for client.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header field) in requests.
func NewClientTLCPFromCert(cp *x509.CertPool, serverNameOverride string) TransportCredentials {
	return NewTLCP(&tlcp.Config{ServerName: serverNameOverride, RootCAs: cp})
}

// NewClientTLCPFromFile constructs TLCP credentials from the input certificate file for client.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header field) in requests.
func NewClientTLCPFromFile(certFile, serverNameOverride string) (TransportCredentials, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return NewTLCP(&tlcp.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}

// NewServerTLCPFromCert constructs TLCP credentials from the input certificate for server.
func NewServerTLCPFromCert(cert *tlcp.Certificate) TransportCredentials {
	return NewTLCP(&tlcp.Config{Certificates: []tlcp.Certificate{*cert}})
}

// NewServerTLCPFromFile constructs TLCP credentials from the input certificate file and key
// file for server.
func NewServerTLCPFromFile(certFile, keyFile string) (TransportCredentials, error) {
	cert, err := tlcp.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewTLCP(&tlcp.Config{Certificates: []tlcp.Certificate{cert}}), nil
}


func cloneTLCPConfig(cfg *tlcp.Config) *tlcp.Config {
	if cfg == nil {
		return &tlcp.Config{}
	}

	return cfg.Clone()
}

