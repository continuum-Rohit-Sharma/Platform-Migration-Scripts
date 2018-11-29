package agent

import "time"

//MailboxMessage is the struct definition of /resources/agent/mailboxMessage
type MailboxMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Path         string    `json:"path"`
	Message      string    `json:"message"`
	MessageID    string
}

//Enum representing status of a mailbox message
const (
	//MailboxMsgStatusPending
	MailboxMsgStatusPending = 0
	//MailboxMsgStatusSeverSent message has been sent down to client
	MailboxMsgStatusSeverSent = 1
	//MailboxMsgStatusAgentProcessedSuccess denotes client has succesfully processed the message
	MailboxMsgStatusAgentProcessedSuccess = 7
	//MailboxMsgStatusAgentProcessedFailure denotes processing failed at client
	MailboxMsgStatusAgentProcessedFailure = 8
)

//MailboxMessageStatus is the struct definition of /resources/agent/mailboxMessage
type MailboxMessageStatus struct {
	StatusCode int `json:"statusCode"`
}
