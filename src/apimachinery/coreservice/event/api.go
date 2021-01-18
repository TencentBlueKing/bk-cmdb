/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (e *event) Subscribe(ctx context.Context, h http.Header, subscription *metadata.Subscription) (*metadata.Subscription, errors.CCErrorCoder) {
	ret := new(metadata.SubscriptionResult)
	subPath := "/create/subscribe"

	err := e.client.Post().
		WithContext(ctx).
		Body(subscription).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil

}

func (e *event) UnSubscribe(ctx context.Context, h http.Header, subscribeID int64) errors.CCErrorCoder {
	ret := new(metadata.SubscriptionResult)
	subPath := "/delete/subscribe/%d"

	err := e.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath, subscribeID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (e *event) UpdateSubscription(ctx context.Context, h http.Header, subscribeID int64, subscription *metadata.Subscription) errors.CCErrorCoder {
	ret := new(metadata.SubscriptionResult)
	subPath := "/update/subscribe/%d"

	err := e.client.Put().
		WithContext(ctx).
		Body(subscription).
		SubResourcef(subPath, subscribeID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (e *event) ListSubscriptions(ctx context.Context, h http.Header, data *metadata.ParamSubscriptionSearch) (*metadata.RspSubscriptionSearch, errors.CCErrorCoder) {
	ret := new(metadata.MultipleSubscriptionResult)
	subPath := "/findmany/subscribe"

	err := e.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}
