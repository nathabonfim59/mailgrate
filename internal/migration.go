package internal

import "fmt"

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

	return nil
}
