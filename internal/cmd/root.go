package cmd

import (
	"os"

	"github.com/spf13/cobra"
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

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
}
