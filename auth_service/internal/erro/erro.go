package erro

import "errors"

var (
	ErrorGetEnvDB                 = errors.New("DB get environment error")
	ErrorNotPost                  = errors.New("Method is not POST")
	ErrorNotGet                   = errors.New("Method is not GET")
	ErrorNotDelete                = errors.New("Method is not DELETE")
	ErrorReadAll                  = errors.New("ReadAll error")
	ErrorUnmarshal                = errors.New("Unmarshal error")
	ErrorMarshal                  = errors.New("Marshal error")
	ErrorNotEmail                 = errors.New("This email format is not supported")
	ErrorUniqueEmail              = errors.New("This email has already been registered")
	ErrorHashPass                 = errors.New("Hash-Password error")
	ErrorInternalServer           = errors.New("Internal Server Error")
	ErrorEmailNotRegister         = errors.New("This email is not registered")
	ErrorFoundUser                = errors.New("Person not found")
	ErrorInvalidPassword          = errors.New("Invalid Password")
	ErrorInvalidSessionID         = errors.New("Session not found")
	ErrorGetSession               = errors.New("Error get session")
	ErrorSetSession               = errors.New("Error set session")
	ErrorGetUserIdSession         = errors.New("UserID not found in session")
	ErrorGetExpirationTimeSession = errors.New("ExpirationTime not found in session")
	ErrorSessionParse             = errors.New("Error parse session data")
	ErrorUnexpectedData           = errors.New("Unexpected data type")
	ErrorStartTransaction         = errors.New("Transaction creation error")
	ErrorCommitTransaction        = errors.New("Transaction commit error")
	ErrorAuthorized               = errors.New("The current session is active")
	ErrorGetUserId                = errors.New("Error getting the UserId from the request context")
	ErrorContextTimeout           = errors.New("The timeout context has expired")
	ErrorSendKafkaMessage         = errors.New("Error Kafka Message")
)
