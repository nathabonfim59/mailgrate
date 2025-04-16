package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nathabonfim59/mailgrate/internal"
	"github.com/spf13/cobra"
)

// Version information set by build flags
var (
	Version   = "dev"
	BuildTime = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mailgrate",
	Short: "A tool for migrating IMAP accounts into Dovecot format",
	Long: `Mailgrate is a utility designed to help migrate email accounts from various IMAP servers
into the Dovecot mail server format. It handles the conversion process while preserving
email metadata, folder structure, and message attributes.

Basic usage:
  mailgrate --source-server imap.example.com --source-user user@example.com --source-pass password \
    --destination-path /var/mail/dovecot/user

With SSL:
  mailgrate --source-server imap.example.com --source-port 993 --use-ssl \
    --source-user user@example.com --source-pass password \
    --destination-path /var/mail/dovecot/user

Include specific folders only:
  mailgrate --source-server imap.example.com --source-user user@example.com \
    --folders "INBOX,Sent,Important" --destination-path /var/mail/dovecot/user

Using a YAML file for multiple users:
  mailgrate --users-file users.yaml --destination-path /var/mail/dovecot`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetVersion returns the current version string
func GetVersion() string {
	return fmt.Sprintf("MailGrate %s (built at %s)", Version, BuildTime)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().BoolP("version", "v", false, "Display version information")

	// Source server configuration
	rootCmd.Flags().String("source-server", "", "Source IMAP server address")
	rootCmd.Flags().Int("source-port", 143, "Source IMAP port (default: 143)")
	rootCmd.Flags().String("source-user", "", "Source IMAP username")
	rootCmd.Flags().String("source-pass", "", "Source IMAP password")
	rootCmd.Flags().Bool("use-ssl", false, "Use SSL/TLS for connection")

	// Migration options
	rootCmd.Flags().String("folders", "", "Comma-separated list of folders to migrate (default: all)")
	rootCmd.Flags().String("destination-path", "", "Path to Dovecot mail directory")
	rootCmd.Flags().Int("concurrent", 5, "Number of concurrent migrations (default: 5)")

	// Multiple users migration
	rootCmd.Flags().String("users-file", "", "Path to YAML file containing multiple users to migrate")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Println(GetVersion())
			return
		}

		// YAML file specified
		if usersFile, _ := cmd.Flags().GetString("users-file"); usersFile != "" {
			fmt.Println("Migrating users from file:", usersFile)
			// err := internal.migrateUsers(usersFile)
			// if err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
			return
		}

		// No YAML file specified, proceed with single user migration
		requiredFlags := []string{"source-server", "source-user", "source-pass"}
		for _, flag := range requiredFlags {
			if value, _ := cmd.Flags().GetString(flag); value == "" {
				fmt.Printf("Error: --%s flag is required\n", flag)
				cmd.Help()
				os.Exit(1)
			}
		}

		sourceServer, _ := cmd.Flags().GetString("source-server")
		sourcePort, _ := cmd.Flags().GetInt("source-port")
		sourceUser, _ := cmd.Flags().GetString("source-user")
		sourcePass, _ := cmd.Flags().GetString("source-pass")
		useSSL, _ := cmd.Flags().GetBool("use-ssl")

		foldersStr, _ := cmd.Flags().GetString("folders")
		var folders []string
		if foldersStr != "" {
			folders = strings.Split(foldersStr, ",")
		}

		destinationPath, _ := cmd.Flags().GetString("destination-path")
		concurrent, _ := cmd.Flags().GetInt("concurrent")

		err := internal.MigrateUser(sourceServer, sourcePort, sourceUser, sourcePass, useSSL, folders, destinationPath, concurrent)

		if err != nil {
			fmt.Println("Error during migration:", err)
			os.Exit(1)
		}

		// err := internal.migrateUser()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		return

		// cmd.Help()
	}
}
