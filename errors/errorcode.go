package errors

type PismoErrorCode string

const (
	PismoErrorCodeInvalidArgument  PismoErrorCode = "INVALID_ARGUMENT"
	PismoErrorCodeNotFound         PismoErrorCode = "NOT_FOUND"
	PismoErrorCodeAlreadyExists    PismoErrorCode = "ALREADY_EXISTS"
	PismoErrorCodePermissionDenied PismoErrorCode = "PERMISSION_DENIED"
	PismoErrorCodeUnauthenticated  PismoErrorCode = "UNAUTHENTICATED"
	PismoErrorCodeInternal         PismoErrorCode = "INTERNAL"
	PismoErrorCodeUnavailable      PismoErrorCode = "UNAVAILABLE"
)
