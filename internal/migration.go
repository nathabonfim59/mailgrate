package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MigrateUser(
	sourceServer string,
	sourcePort int,
	sourceUser string,
	sourcePass string,
	useSSL bool,
	folders []string,
	destinationPath string,
	concurrent int,
	provider Provider,
) error {
	// Connect to the source IMAP server
	server, err := connect(sourceServer, sourcePort, sourceUser, sourcePass, useSSL)

	if err != nil {
		return fmt.Errorf("failed to connect to source server: %w", err)
	}

	fmt.Println("Connected to source server: ", server.Client.State().String())

	inboxes, err := server.ListMailboxes()
	if err != nil {
		return fmt.Errorf("failed to list mailboxes: %w", err)
	}

	err = os.MkdirAll(destinationPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	// Filter mailboxes if specific folders were requested
	var filteredInboxes []*struct {
		Mailbox    string
		Attributes []string
	}
	if len(folders) > 0 {
		mapping := GetProviderMailboxMapping(provider)

		for _, inbox := range inboxes {
			// Check if this mailbox is in the requested folders list
			for _, folder := range folders {
				// Map standard folder names to provider-specific names
				providerFolder := folder
				switch strings.ToLower(folder) {
				case "inbox":
					providerFolder = mapping.Inbox
				case "drafts":
					providerFolder = mapping.Drafts
				case "sent":
					providerFolder = mapping.Sent
				case "spam":
					providerFolder = mapping.Spam
				case "trash":
					providerFolder = mapping.Trash
				}

				// Add the mailbox if it matches
				if inbox.Mailbox == providerFolder {
					filteredInboxes = append(filteredInboxes, &inbox)
					break
				}
			}
		}
	} else {
		// If no specific folders requested, use all mailboxes
		filteredInboxes = inboxes
	}

	for _, inbox := range filteredInboxes {
		fmt.Println("Saving emails from mailbox:", inbox.Mailbox)

		server.SelectMailbox(inbox.Mailbox)

		// Get standardized mailbox name based on provider
		standardMailboxName := GetStandardMailboxName(provider, inbox.Mailbox)

		// Create the mailbox directory using standardized name
		mailboxPath := filepath.Join(destinationPath, standardMailboxName)
		err = os.MkdirAll(mailboxPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create mailbox directory: %w", err)
		}

		// Download all messages in the mailbox with the Dovecot format
		err = server.DownloadAllMessages(mailboxPath, DovecotFormat)
		if err != nil {
			return fmt.Errorf("failed to download messages from %s: %w", inbox.Mailbox, err)
		}
	}

	return nil
}
