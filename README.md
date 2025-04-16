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
- `--help`: Show help information

## License

MIT License - Copyright (c) Nathanael Bonfim

