package internal

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// Message represents an email message with its details
type Message struct {
	UID           uint32
	SeqNum        uint32
	Subject       string
	From          string
	To            string
	Date          time.Time
	Size          int64
	Flags         []imap.Flag
	MessageID     string
	HasAttachment bool
}

// ListMessages retrieves all messages in the selected mailbox
func (s *IMAPServer) ListMessages() ([]*Message, error) {
	// Verify that a mailbox has been selected
	if s.Client.State() != imap.ConnStateSelected {
		return nil, errors.New("no mailbox selected")
	}

	seqSet := imap.SeqSet{}
	seqSet.AddRange(1, 0) // 0 is special and means the last message

	// Configure fetch options to retrieve message details
	fetchOptions := &imap.FetchOptions{
		UID:           true,
		Envelope:      true,
		Flags:         true,
		RFC822Size:    true,
		BodyStructure: &imap.FetchItemBodyStructure{},
	}

	// Execute fetch command
	fetchCmd := s.Client.Fetch(seqSet, fetchOptions)
	defer fetchCmd.Close()

	// Collect the message data
	msgBuffers, err := fetchCmd.Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}

	messages := make([]*Message, 0, len(msgBuffers))
	for _, msg := range msgBuffers {
		message := &Message{
			UID:       uint32(msg.UID),
			SeqNum:    msg.SeqNum,
			Flags:     msg.Flags,
			Size:      msg.RFC822Size,
			Date:      msg.Envelope.Date,
			MessageID: msg.Envelope.MessageID,
		}

		// Extract subject
		if msg.Envelope.Subject != "" {
			message.Subject = msg.Envelope.Subject
		}

		// Extract from address
		if len(msg.Envelope.From) > 0 {
			addr := msg.Envelope.From[0]
			if addr.Mailbox != "" && addr.Host != "" {
				if addr.Name != "" {
					message.From = fmt.Sprintf("%s <%s@%s>", addr.Name, addr.Mailbox, addr.Host)
				} else {
					message.From = fmt.Sprintf("%s@%s", addr.Mailbox, addr.Host)
				}
			}
		}

		// Extract to address
		if len(msg.Envelope.To) > 0 {
			addr := msg.Envelope.To[0]
			if addr.Mailbox != "" && addr.Host != "" {
				if addr.Name != "" {
					message.To = fmt.Sprintf("%s <%s@%s>", addr.Name, addr.Mailbox, addr.Host)
				} else {
					message.To = fmt.Sprintf("%s@%s", addr.Mailbox, addr.Host)
				}
			}
		}

		// Check if the message has attachments by examining the body structure
		if msg.BodyStructure != nil {
			message.HasAttachment = hasAttachment(msg.BodyStructure)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// hasAttachment checks if a message has attachments by examining its body structure
func hasAttachment(bs imap.BodyStructure) bool {
	hasAttach := false

	// Walk through the body structure parts
	bs.Walk(func(path []int, part imap.BodyStructure) bool {
		// Check for specific content dispositions that indicate attachments
		if part.Disposition() != nil && part.Disposition().Value == "attachment" {
			hasAttach = true
			return false // Stop traversing
		}

		// Check content type for typical attachment types
		mediaType := part.MediaType()
		if strings.HasPrefix(mediaType, "application/") &&
			!strings.Contains(mediaType, "application/pgp-signature") {
			hasAttach = true
			return false
		}

		return true // Continue traversing
	})

	return hasAttach
}

// DownloadMessage downloads a specific message as an EML file
func (s *IMAPServer) DownloadMessage(uid uint32, outputPath string) error {
	// Verify that a mailbox has been selected
	if s.Client.State() != imap.ConnStateSelected {
		return errors.New("no mailbox selected")
	}

	// Create a UID set for the specific message
	uidSet := imap.UIDSet{}
	uidSet.AddNum(imap.UID(uid))

	// Configure fetch options to retrieve the entire message
	fetchOptions := &imap.FetchOptions{
		UID: true,
		BodySection: []*imap.FetchItemBodySection{
			// Fetch the entire message including headers and body
			{},
		},
	}

	// Execute fetch command
	fetchCmd := s.Client.Fetch(uidSet, fetchOptions)
	defer fetchCmd.Close()

	// Collect the message data
	msgBuffers, err := fetchCmd.Collect()
	if err != nil {
		return fmt.Errorf("failed to fetch message: %v", err)
	}

	if len(msgBuffers) == 0 {
		return fmt.Errorf("no message found with UID %d", uid)
	}

	// Get the first (and should be only) message
	msg := msgBuffers[0]

	// Check if we have the body section data
	if len(msg.BodySection) == 0 {
		return errors.New("message body not received")
	}

	// Get the message body data - use empty FetchItemBodySection{} and FindBodySection
	// to retrieve the entire message content
	emptySection := &imap.FetchItemBodySection{}
	bodyData := msg.FindBodySection(emptySection)

	if bodyData == nil || len(bodyData) == 0 {
		return errors.New("message body not found")
	}

	// Create a reader for the body data
	bodyReader := bytes.NewReader(bodyData)

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Write the message data to the file
	_, err = io.Copy(file, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to write message to file: %v", err)
	}

	return nil
}

// DownloadAllMessages downloads all messages in the selected mailbox as EML files
func (s *IMAPServer) DownloadAllMessages(outputDir string) error {
	// First get the list of all messages
	messages, err := s.ListMessages()
	if err != nil {
		return fmt.Errorf("failed to list messages: %v", err)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Download each message
	for _, msg := range messages {
		// Create a filename based on the UID and subject
		// Clean the subject to make it safe for filenames
		subject := msg.Subject
		if subject == "" {
			subject = "no_subject"
		}
		// Remove characters that are not safe for filenames
		subject = strings.Map(func(r rune) rune {
			if strings.ContainsRune(`<>:"/\|?*`, r) {
				return '_'
			}
			return r
		}, subject)

		// Limit the subject length for the filename
		if len(subject) > 50 {
			subject = subject[:50]
		}

		// Create the filename
		filename := fmt.Sprintf("%d-%s.eml", msg.UID, subject)
		outputPath := filepath.Join(outputDir, filename)

		// Download the message
		if err := s.DownloadMessage(msg.UID, outputPath); err != nil {
			return fmt.Errorf("failed to download message %d: %v", msg.UID, err)
		}
	}

	return nil
}
