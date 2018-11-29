package protocol

// HeaderKey is protocol header
type HeaderKey string

//ResponseStatus stataus code for responses from the server
type ResponseStatus int

// Headers constant values
const (
	//HdrConstForceSend is a constant for header HdrForceSend
	HdrConstForceSend string = "true"
	//HdrConstPluginDataPersist is a constant for header HdrPluginDataPersist
	HdrConstPluginDataPersist string = "true"
)

const (
	// HdrUserAgent describes client making protocol request
	HdrUserAgent HeaderKey = "User-Agent"

	// HdrContentType describes type of content in request or response
	HdrContentType HeaderKey = "Content-Type"

	//HdrPluginDataPersist describes whether to persist plugin data if server is offline
	HdrPluginDataPersist HeaderKey = "Continuum-Plugin-Persist-Data"

	// HdrForceSend describes whether to try sending data even if server is offline
	HdrForceSend HeaderKey = "Continuum-Plugin-Force-Send"

	// HdrErrorCode is for top level error code for a failed request
	HdrErrorCode HeaderKey = "Continuum-Plugin-Error-Code"

	//HdrContentMD5 is MD5 hash key
	HdrContentMD5 HeaderKey = "Content-MD5"

	//HdrContentWebhook indicates message contains webhook
	HdrContentWebhook HeaderKey = "Content-Webhook"

	//HdrDataCompressionType indicates message contains data compression type
	HdrDataCompressionType = "Data-Compression-Type"

	//HdrAcceptEncoding indicates message contains Accept Encoding
	HdrAcceptEncoding = "Accept-Encoding"

	//Ok Response status for sucessful response
	Ok ResponseStatus = 200

	//StatusCreated Response status for Created response
	StatusCreated ResponseStatus = 201

	//StatusCodeInternalError status for internal exception in executing the task
	StatusCodeInternalError ResponseStatus = 500

	//StatusCodeBadRequest error status for bad request
	StatusCodeBadRequest ResponseStatus = 400

	//PathNotFound Error status for incorrect Plugin Path
	PathNotFound ResponseStatus = 404

	//HdrBrokerPath describes Broker URL where the data would be posted
	HdrBrokerPath HeaderKey = "Continuum-Plugin-Broker-Path"

	//HdrCommunicationPath describes Broker URL where the data would be posted
	HdrCommunicationPath HeaderKey = "Continuum-Plugin-Communication-Path"

	//HdrTaskInput describes Task Input where for execution
	HdrTaskInput HeaderKey = "Continuum-Plugin-Task-Input"

	//HdrMessageType describes Message Type to process mailbox message at plugin
	HdrMessageType HeaderKey = "Continuum-Plugin-Message-Type"

	//HdrTransactionID describes RequestID/TransactionID/CorreleationID to track data accross servers and processes.
	HdrTransactionID HeaderKey = "X-Request-Id"

	//HdrHTTPSecure This is temporary Key used for heartbeat, would be removed once the heartbeat changes are done in communication service
	HdrHTTPSecure HeaderKey = "Continuum-HTTP-Secure"

	//HdrAgentOS : This is a header key to pass Agent OS; as a part of any request from Agent
	HdrAgentOS string = "Continuum-Agent-OS"

	//HdrPluginTimeout :  This is a header key to pass timeout to the respective plugin
	HdrPluginTimeout HeaderKey = "Plugin-Execution-Timeout"

	//HdrResourcePath describes Resource URL for which data is to be fetched
	HdrResourcePath HeaderKey = "Continuum-Plugin-Resource-Path"
)

// Headers is a map for Request Response structures
type Headers map[HeaderKey][]string

// SetKeyValue sets a key to a value (overwriting if it exists)
func (h Headers) SetKeyValue(key HeaderKey, value string) {
	h.SetKeyValues(key, []string{value})
}

// SetKeyValues sets a key to values (overwrting if it exists)
func (h Headers) SetKeyValues(key HeaderKey, values []string) {
	h[key] = values
}

// GetKeyValue returns the value for a given key
func (h Headers) GetKeyValue(key HeaderKey) (value string) {
	values := h[key]
	if len(values) > 0 {
		value = values[0]
	}
	return
}

// GetKeyValues returns values array for given key
func (h Headers) GetKeyValues(key HeaderKey) (values []string) {
	return h[key]
}

// Parameters is a map of request parameters
type Parameters map[string][]string
