package agent

//WebhookScheduleTask is a struct defining the actual webhookScheduleTask message
type WebhookScheduleTask struct {
	/** subType contains the value of the webhook schedule task */
	SubType string `json:"subType"`
	/** action hold the value of the type of operation need to be perform */
	Action string `json:"action"`
	/**webhook contains the URL */
	Webhook string `json:"webhook"`
	/** Timeout value for webhook schedule task*/
	ScheduleTimeoutInSeconds int `json:"scheduleTimeoutInSeconds"`
	/** DataCompressionType  like gzip,deflate*/
	DataCompressionType string `json:"dataCompressionType"`
}
