package internal

import (
	"fmt"
	"os"
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

	for _, inbox := range inboxes {
		fmt.Println("Saving emails from Inbox: ", inbox.Mailbox)

		server.SelectMailbox(inbox.Mailbox)

		// Remove the 'INBOX' prefix from the mailbox name (not used in dovecot format)
		mailboxName := strings.TrimPrefix(inbox.Mailbox, "INBOX.")
		if mailboxName == "" {
			mailboxName = "INBOX"
		}

		// Create the mailbox directory
		mailboxPath := fmt.Sprintf("%s/%s", destinationPath, mailboxName)
		err = os.MkdirAll(mailboxPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create mailbox directory: %w", err)
		}

		// Download all messages in the mailbox with the Dovecot format
		server.DownloadAllMessages(mailboxPath, DovecotFormat)
	}

	return nil
}
