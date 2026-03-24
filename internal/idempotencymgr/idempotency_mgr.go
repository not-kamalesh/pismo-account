package idempotencymgr

// This can be extended with the redis(or some other datastore) as a store for idempotency entries
// when we want to scale to multiple servers

//go:generate mockery --name=IdempotencyMgr --output=. --outpkg=idempotencymgr --filename=mock_idempotencymgr.go --structname=MockIdempotencyMgr
type IdempotencyMgr interface {
	Execute(idempotencyKey string, reqHash string, handler func() (interface{}, error)) (interface{}, error)
}
