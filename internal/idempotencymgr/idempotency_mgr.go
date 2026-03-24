package idempotencymgr

// This can be extended with the redis(or some other datastore) as a store for idempotency entries
// when we want to scale to multiple servers

//go:generate mockery --name=IdempotencyMgr --output=. --outpkg=idempotencymgr --filename=mock_idempotencymgr.go --structname=MockIdempotencyMgr

// IdempotencyMgr provides an interface to process requests in a idempotent manner, meaning same request should get same response
// will be used in the transaction apis
type IdempotencyMgr interface {
	Execute(idempotencyKey string, reqHash string, handler func() (interface{}, error)) (interface{}, error)
}
