package internal

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// IMAPServer represents a connected IMAP server session
type IMAPServer struct {
	Client *imapclient.Client
}

// connect establishes a connection to an IMAP server based on command flags
func connect(server string, port int, username string, password string, useSSL bool) (*IMAPServer, error) {
	// Validate required connection parameters
	if server == "" {
		return nil, errors.New("source server is required")
	}

	if username == "" {
		return nil, errors.New("source username is required")
	}

	if password == "" {
		return nil, errors.New("source password is required")
	}

	// Validate required connection parameters
	if server == "" {
		return nil, errors.New("source server is required")
	}

	if username == "" {
		return nil, errors.New("source username is required")
	}

	if password == "" {
		return nil, errors.New("source password is required")
	}

	// Prepare connection options
	options := &imapclient.Options{
		// For production use, you might want to set up proper TLS configuration
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // Consider setting to false in production
		},
	}

	// Connect to the server based on selected method
	var client *imapclient.Client
	var err error
	address := fmt.Sprintf("%s:%d", server, port)

	if useSSL {
		// Use TLS connection
		client, err = imapclient.DialTLS(address, options)
	} else {
		// Use insecure connection with option to upgrade
		client, err = imapclient.DialInsecure(address, options)

		// If server supports STARTTLS, upgrade the connection
		if err == nil && client.Caps().Has(imap.CapStartTLS) {
			// Close the current connection and create a new one with TLS
			client.Close()
			client, err = imapclient.DialStartTLS(address, options)
			if err != nil {
				return nil, fmt.Errorf("failed to establish secure connection using STARTTLS: %v", err)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %v", err)
	}

	// Login to server
	if err := client.Login(username, password).Wait(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	return &IMAPServer{
		Client: client,
	}, nil
}

// Close closes the IMAP connection
func (s *IMAPServer) Close() error {
	if s.Client != nil {
		return s.Client.Close()
	}
	return nil
}

// ListMailboxes lists all mailboxes on the server
func (s *IMAPServer) ListMailboxes() ([]*imap.ListData, error) {
	cmd := s.Client.List("", "*", nil)
	defer cmd.Close()

	return cmd.Collect()
}

// SelectMailbox selects a mailbox for reading messages
func (s *IMAPServer) SelectMailbox(mailbox string) (*imap.SelectData, error) {
	return s.Client.Select(mailbox, nil).Wait()
}
