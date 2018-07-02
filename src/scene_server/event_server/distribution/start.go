package distribution

import (
	redis "gopkg.in/redis.v5"

	"configcenter/src/scene_server/event_server/identifier"
	"configcenter/src/storage"
)

func Start(cache *redis.Client, db storage.DI) error {
	chErr := make(chan error)

	eh := &EventHandler{cache: cache}
	go func() {
		chErr <- eh.StartHandleInsts()
	}()

	dh := &DistHandler{cache: cache}
	go func() {
		chErr <- dh.StartDistribute()
	}()

	ih := identifier.NewIdentifierHandler(cache, db)
	go func() {
		chErr <- ih.StartHandleInsts()
	}()

	return <-chErr
}

type EventHandler struct {
	cache *redis.Client
}

type DistHandler struct {
	cache *redis.Client
}
