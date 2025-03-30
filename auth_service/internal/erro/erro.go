package erro

import "errors"

var (
	ErrorGetEnvDB                 = errors.New("DB get environment error")
	ErrorNotPost                  = errors.New("Method is not POST")
	ErrorReadAll                  = errors.New("ReadAll error")
	ErrorUnmarshal                = errors.New("Unmarshal error")
	ErrorNotEmail                 = errors.New("This email format is not supported")
	ErrorUniqueEmail              = errors.New("This email has already been registered")
	ErrorHashPass                 = errors.New("Hash-Password error")
	ErrorInternalServer           = errors.New("Internal Server Error")
	ErrorEmailNotRegister         = errors.New("This email is not registered")
	ErrorInvalidPassword          = errors.New("Invalid Password")
	ErrorInvalidSessionID         = errors.New("Session not found")
	ErrorGetSession               = errors.New("Error get session")
	ErrorSetSession               = errors.New("Error set session")
	ErrorGetUserIdSession         = errors.New("UserID not found in session")
	ErrorGetExpirationTimeSession = errors.New("ExpirationTime not found in session")
	ErrorSessionParse             = errors.New("Error parse session data")
	ErrorUnexpectedData           = errors.New("Unexpected data type")
)
