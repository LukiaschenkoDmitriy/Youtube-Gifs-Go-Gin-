package errorcode

type ErrorCode int

const (
	AlreadyAuthorized      = ErrorCode(-1)
	NotAuthorized          = ErrorCode(iota)
	EntityNotFound         = ErrorCode(1)
	ClientError            = ErrorCode(2)
	EntityCreationFailed   = ErrorCode(3)
	WrongAcceptHeader      = ErrorCode(4)
	WrongJSONDataStructure = ErrorCode(5)
	DataBaseError          = ErrorCode(6)
)
