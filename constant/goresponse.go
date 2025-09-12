package constant

const (
	// Protocol
	ProtocolWebApi = "web-api"
	ProtocolGrpc   = "grpc"

	EnLanguage = "en"
	IdLanguage = "id"

	// Language
	DefaultLanguage = EnLanguage
)

const (
	IsResponseSuccess   = "success"
	IsResponseCreated   = "created"
	IsResponseUpdated   = "updated"
	IsResponseDeleted   = "deleted"
	IsResponseRetrieved = "retrieved"

	ErrorInternalServer     = "internal_server_error"
	ErrorConnectionRefused  = "connection_refused"
	ErrorTooManyConnections = "too_many_connections"
	ErrorConnectionTimeout  = "connection_timeout"
	ErrorDnsError           = "dns_error"
	ErrorAuthFailed         = "auth_failed"
	ErrorAccessDenied       = "access_denied"
	ErrorDriverError        = "driver_error"
	ErrorSslTlsError        = "ssl_tls_error"
)
