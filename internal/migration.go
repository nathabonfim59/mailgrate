package internal

import (
	"fmt"

	"github.com/sanity-io/litter"
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

	for _, inbox := range inboxes {
		fmt.Println("Inbox: ", inbox)

		server.SelectMailbox(inbox.Mailbox)

		// Fetch messages from the selected mailbox

		litter.Dump(inbox)
	}

	return nil
}
