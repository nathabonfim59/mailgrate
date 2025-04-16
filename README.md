# Mailgrate

A tool for migrating IMAP accounts into Dovecot format.

## Description

Mailgrate is a utility designed to help migrate email accounts from various IMAP servers into the Dovecot mail server format. It handles the conversion process while preserving email metadata, folder structure, and message attributes.

## Installation

```bash
# Clone the repository
git clone https://github.com/nathanaelbonfim/mailgrate.git

# Build the binary
make build
```

## Usage

```bash
# Basic usage
mailgrate --source-server imap.example.com --source-user user@example.com --source-pass password \
  --destination-path /var/mail/dovecot/user

# With SSL
mailgrate --source-server imap.example.com --source-port 993 --use-ssl \
  --source-user user@example.com --source-pass password \
  --destination-path /var/mail/dovecot/user

# Include specific folders only
mailgrate --source-server imap.example.com --source-user user@example.com \
  --folders "INBOX,Sent,Important" --destination-path /var/mail/dovecot/user

# Using a YAML file for multiple users
mailgrate --users-file users.yaml --destination-path /var/mail/dovecot
```

## Options

- `--source-server`: Source IMAP server address
- `--source-port`: Source IMAP port (default: 143)
- `--source-user`: Source IMAP username
- `--source-pass`: Source IMAP password
- `--use-ssl`: Use SSL/TLS for connection
- `--folders`: Comma-separated list of folders to migrate (default: all)
- `--destination-path`: Path to Dovecot mail directory
- `--concurrent`: Number of concurrent migrations (default: 5)
- `--users-file`: Path to YAML file containing multiple users to migrate
- `--help`: Show help information

## Multiple User Migration

You can specify multiple users in a YAML file for batch migration. Create a YAML file with the following structure:

```yaml
# Example users.yaml file
hosts:
  - server: imap.example.com
    port: 143
    use_ssl: false
    users:
      - email: user1@example.com
        password: password1
        folders: ["INBOX", "Sent", "Important"]
      - email: user2@example.com
        password: password2
        
  - server: imap.another.com
    port: 993
    use_ssl: true
    users:
      - email: user3@another.com
        password: password3
        folders: ["INBOX", "Archive"]
      - email: user4@another.com
        password: password4
```

Then run the migration using:

```bash
mailgrate --users-file users.yaml --destination-path /var/mail/dovecot --concurrent 10
```

This will process all users defined in the YAML file, maintaining the folder structure under the destination path.

## License

MIT License - Copyright (c) Nathanael Bonfim
