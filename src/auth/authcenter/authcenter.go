package authcenter

import "configcenter/src/auth"

// authCenter means BlueKing's authorize center,
// which is also a open source product.
type authCenter struct {
}

func (ac *authCenter) Register(r *auth.ResourceAttribute) (requestID string, err error) {
	panic("implement me")
}

func (ac *authCenter) Deregister(r *auth.ResourceAttribute) (requestID string, err error) {
	panic("implement me")
}

func (ac *authCenter) Update(r *auth.ResourceAttribute) (requestID string, err error) {
	panic("implement me")
}

func (ac *authCenter) Get() error {
	panic("implement me")
}
