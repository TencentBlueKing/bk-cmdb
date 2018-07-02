package distribution

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/event_server/types"
)

func (dh *DistHandler) SendCallback(receiver *metadata.Subscription, event string) (err error) {
	dh.cache.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "total", 1)

	body := bytes.NewBufferString(event)
	req, err := http.NewRequest("POST", receiver.CallbackURL, body)
	if err != nil {
		dh.cache.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
		return fmt.Errorf("event distribute fail, build request error: %v, date=[%s]", err, event)
	}
	var duration time.Duration
	if receiver.TimeOut == 0 {
		duration = timeout
	} else {
		duration = receiver.GetTimeout()
	}
	resp, err := httpCli.DoWithTimeout(duration, req)
	if err != nil {
		dh.cache.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
		return fmt.Errorf("event distribute fail, send request error: %v, date=[%s]", err, event)
	}
	defer resp.Body.Close()
	respdata, _ := ioutil.ReadAll(resp.Body)
	if receiver.ConfirmMode == metadata.ConfirmmodeHttpstatus {
		if strconv.Itoa(resp.StatusCode) != receiver.ConfirmPattern {
			dh.cache.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
			return fmt.Errorf("event distribute fail, received response %s, date=[%s]", respdata, event)
		}
	} else if receiver.ConfirmMode == metadata.ConfirmmodeRegular {
		pattern, err := regexp.Compile(receiver.ConfirmPattern)
		if err != nil {
			return fmt.Errorf("event distribute fail, build regexp error: %v", err)
		}
		if !pattern.Match(respdata) {
			dh.cache.HIncrBy(types.EventCacheDistCallBackCountPrefix+fmt.Sprint(receiver.SubscriptionID), "failue", 1)
			return fmt.Errorf("event distribute fail, received response %s, date=[%s]", respdata, event)
		}
		return nil
	}

	return
}

var httpCli = httpclient.NewHttpClient()
