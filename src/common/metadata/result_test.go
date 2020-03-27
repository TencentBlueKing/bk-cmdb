package metadata

import (
	"testing"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func TestResponse(t *testing.T) {

	err := errors.New(9999999, "test-msg")

	respPtr := &metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			Code:   err.GetCode(),
			ErrMsg: err.Error(),
		},
	}

	ccErr := respPtr.CCError()
	if ccErr == nil {
		t.Errorf("not error")
		return
	}

	if err.GetCode() != ccErr.GetCode() ||
		err.Error() != ccErr.Error() {
		t.Errorf("code info, code:%v, error msg:%s", ccErr.GetCode(), ccErr.Error())
		return
	}

	resp := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			Code:   err.GetCode(),
			ErrMsg: err.Error(),
		},
	}

	ccErr = resp.CCError()
	if ccErr == nil {
		t.Errorf("not error")
		return
	}

	if err.GetCode() != ccErr.GetCode() ||
		err.Error() != ccErr.Error() {
		t.Errorf("code info, code:%v, error msg:%s", ccErr.GetCode(), ccErr.Error())
		return
	}

	respSucc := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
			ErrMsg: "",
		},
	}
	ccErr = respSucc.CCError()
	if ccErr != nil {
		t.Errorf("have error")
		return
	}

}
