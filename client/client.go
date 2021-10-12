package client

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type Status int8

type client struct {
	host   string
	logger *zerolog.Logger
}

type Client interface {
	GetExpiry() (Status, time.Time)
}

const (
	Ok Status = iota
	ErrConnection
	ErrTimeout
	ErrCertInvalid
	ErrCertExpired
)

func New(hostname string, zerolog *zerolog.Logger) Client {
	return &client{
		host:   hostname,
		logger: zerolog,
	}
}

func (c *client) GetExpiry() (Status, time.Time) {
	protocol := "tcp"
	hostPort := c.host + ":443"
	c.logInfo("Establishing TLS connection to " + protocol + "://" + hostPort)

	tlsConn, err := tls.Dial(protocol, hostPort, nil)
	if err != nil {
		c.logWarn(err)
		return ErrConnection, time.Now()
	}
	defer tlsConn.Close()

	if err = tlsConn.VerifyHostname(c.host); err != nil {
		c.logWarn(err)
		return ErrCertInvalid, time.Now()
	}

	expiry := tlsConn.ConnectionState().PeerCertificates[0].NotAfter
	if time.Until(expiry) <= 0 {
		c.logWarn(errors.New("certificate has been expired"))
		return ErrCertExpired, expiry
	}

	c.logInfo("Operation succeeded")
	return Ok, expiry
}

func (c *client) logWarn(err error) {
	if c.logger != nil {
		c.logger.Warn().Msgf("%s: %s", c.host, err.Error())
	}
}

func (c *client) logInfo(info string) {
	if c.logger != nil {
		c.logger.Info().Msgf("%s: %s", c.host, info)
	}
}
