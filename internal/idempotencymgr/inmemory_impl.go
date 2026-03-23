package idempotencymgr

import (
	"log/slog"
	"sync"
	"time"

	"github.com/not-kamalesh/pismo-account/errors"
)

type IdempotencyEntry struct {
	RequestHash string
	Response    interface{}
	Err         error
	Status      string
	Done        chan struct{}
	CreatedAt   time.Time
}

type InMemIdempotencyMgr struct {
	mu    sync.Mutex
	store map[string]*IdempotencyEntry
	ttl   time.Duration
}

func NewInMemIdempotencyMgr() IdempotencyMgr {
	iMgr := &InMemIdempotencyMgr{
		store: make(map[string]*IdempotencyEntry),
		ttl:   5 * time.Minute,
	}
	iMgr.startCleanup()
	return iMgr
}

func (i *InMemIdempotencyMgr) Execute(idempotencyKey string, reqHash string, handler func() (interface{}, error)) (interface{}, error) {
	// Take a lock
	i.mu.Lock()

	// Check if entry exists for the idempotencyKey
	if entry, ok := i.store[idempotencyKey]; ok {
		// if the request has does not match, return a conflict error
		slog.Info("both hashes", "currHash", reqHash, "storedHash", entry.RequestHash)
		if entry.RequestHash != reqHash {
			i.mu.Unlock()
			return nil, errors.ErrTransactionConflict
		}

		// if the entry is in processing state, wait for the other thread to complete
		// and then respond with the stored status (provides better user experience)
		if entry.Status == "PROCESSING" {
			done := entry.Done
			i.mu.Unlock()
			<-done // wait for done channel to close
			return entry.Response, entry.Err
		}

		// If already in Completed state, return the stored resp and err
		resp, err := entry.Response, entry.Err
		i.mu.Unlock()
		return resp, err
	}

	// Create new entry in the idempotency store
	entry := &IdempotencyEntry{
		RequestHash: reqHash,
		Status:      "PROCESSING",
		Done:        make(chan struct{}),
		CreatedAt:   time.Now(),
	}

	i.store[idempotencyKey] = entry
	i.mu.Unlock()

	var resp any
	var err error

	// Store the resp and err before returning
	// Also close the done channel, so other threads waiting can proceed
	defer func() {
		i.mu.Lock()
		entry.Response = resp
		entry.Err = err
		entry.Status = "COMPLETED"

		close(entry.Done)
		i.mu.Unlock()
	}()

	resp, err = handler()

	return resp, err
}

func (i *InMemIdempotencyMgr) startCleanup() {
	go func() {
		ticker := time.NewTicker(i.ttl)

		for range ticker.C {
			i.mu.Lock()
			for k, v := range i.store {
				if time.Since(v.CreatedAt) > i.ttl {
					delete(i.store, k)
				}
			}
			i.mu.Unlock()
		}
	}()
}
