package util

var (
	NoSessionError                  = "cannot find session in context"
	NoUserIDError                   = "cannot find user_id in session"
	NoSessionIDError                = "cannot find session_id in session"
	NoCSRFError                     = "cannot find csrf_token in context"
	NoCartError                     = "cannot find cart in context"
	HTTPErrorAssertErr              = "failed to assert type: echo.HTTPError"
	NoMenuItem                      = "menu not found. please order in the menu"
	InvalidNotificationSignatureKey = "invalid transaction signature_key"
	NoURLFound                      = "cannot found request URL"
	InvalidRequestOperation         = "invalid request operation"
)
