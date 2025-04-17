package internal

// Provider represents an email service provider
type Provider string

const (
	// DefaultProvider is used when no specific provider is specified
	DefaultProvider Provider = "default"

	// LocawebProvider represents Locaweb mail service
	LocawebProvider Provider = "locaweb"

	// GmailProvider represents Gmail service
	GmailProvider Provider = "gmail"

	// OutlookProvider represents Outlook/Office365 service
	OutlookProvider Provider = "outlook"
)

// MailboxMapping holds the mapping between provider-specific mailbox names and standard names
type MailboxMapping struct {
	Inbox  string
	Drafts string
	Sent   string
	Spam   string
	Trash  string
	// Add more standard folders as needed
}

// GetProviderMailboxMapping returns the standard mailbox mapping for a given provider
func GetProviderMailboxMapping(provider Provider) MailboxMapping {
	switch provider {
	case LocawebProvider:
		return MailboxMapping{
			Inbox:  "INBOX",
			Drafts: "INBOX.rascunho",
			Sent:   "INBOX.enviadas",
			Spam:   "INBOX.Mala_Direta",
			Trash:  "INBOX.lixo",
		}
	case GmailProvider:
		return MailboxMapping{
			Inbox:  "INBOX",
			Drafts: "[Gmail]/Drafts",
			Sent:   "[Gmail]/Sent Mail",
			Spam:   "[Gmail]/Spam",
			Trash:  "[Gmail]/Trash",
		}
	case OutlookProvider:
		return MailboxMapping{
			Inbox:  "INBOX",
			Drafts: "Drafts",
			Sent:   "Sent Items",
			Spam:   "Junk Email",
			Trash:  "Deleted Items",
		}
	default:
		return MailboxMapping{
			Inbox:  "INBOX",
			Drafts: "Drafts",
			Sent:   "Sent",
			Spam:   "Spam",
			Trash:  "Trash",
		}
	}
}

// GetStandardMailboxName returns the standardized mailbox name for a given provider-specific name
func GetStandardMailboxName(provider Provider, mailboxName string) string {
	// Get the mapping for the provider
	mapping := GetProviderMailboxMapping(provider)

	// Check if this mailbox name matches any standard folder
	switch mailboxName {
	case mapping.Inbox:
		return "INBOX"
	case mapping.Drafts:
		return "Drafts"
	case mapping.Sent:
		return "Sent"
	case mapping.Spam:
		return "Spam"
	case mapping.Trash:
		return "Trash"
	default:
		// If no specific mapping, return the original name
		// but remove INBOX prefix if present (common in many providers)
		return mailboxName
	}
}
