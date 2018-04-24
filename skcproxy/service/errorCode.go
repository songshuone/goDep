package service

const (
	ErrInvalidRequest  =iota
	ErrNotFoundProductId
	ErrActiveAlreadyEnd
	ErrActiveNotStart
	ErrActiveSaleOut
	ErrUserCheckAuthFailed
	ErrUserServiceBusy
	ErrProcessTimeout
	ErrClientClosed
)
