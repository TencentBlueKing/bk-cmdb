package distribution

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"configcenter/src/common/core/cc/api"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage"
)

func initTester() {
	config := map[string]string{
		storage.DI_REDIS + ".host":     "127.0.0.1",
		storage.DI_REDIS + ".port":     "6379",
		storage.DI_REDIS + ".usr":      "cc",
		storage.DI_REDIS + ".pwd":      "cc",
		storage.DI_REDIS + ".database": "0",

		storage.DI_MONGO + ".host":     "127.0.0.1",
		storage.DI_MONGO + ".port":     "27017",
		storage.DI_MONGO + ".usr":      "cc",
		storage.DI_MONGO + ".pwd":      "cc",
		storage.DI_MONGO + ".database": "0",
	}
	a := api.NewAPIResource()
	a.GetDataCli(config, storage.DI_REDIS)
	a.GetDataCli(config, storage.DI_MONGO)
}

func TestSendCallback(t *testing.T) {
	initTester()
	f := func(http.ResponseWriter, *http.Request) {}
	s := httptest.NewServer(http.HandlerFunc(f))
	defer s.Close()
	var receiver = &types.Subscription{
		CallbackURL:    s.URL,
		ConfirmMode:    types.ConfirmmodeHttpstatus,
		ConfirmPattern: "200",
		TimeOut:        10,
	}
	if err := SendCallback(receiver, "test message"); err != nil {
		t.Fail()
	}

}
