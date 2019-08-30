/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

// Auth used to implement permission resource registration
type Auth struct {
}

//
// // CheckAuthWithObject
// func (a Auth) CheckAuthWithObject(ctx types.ContextParams, targetObjects []model.Object) (bool, error) {
//
// 	authAttr := &auth.Attribute{
// 		User: auth.UserInfo{
// 			UserName:   ctx.User,
// 			SupplierID: ctx.SupplierAccount,
// 		},
// 	}
//
// 	for _, obj := range targetObjects {
// 		authAttr.Resources = append(authAttr.Resources, auth.Resource{
// 			Name: obj.Object().ObjectName,
// 		})
// 	}
//
// 	descision, reason, err := ctx.AuthAPI.Authorize(authAttr)
//
// 	// TODO: need to be implemented, unknown how to use the desision and reason value
// 	_ = descision
// 	_ = reason
//
// 	panic("need to implemented")
//
// 	return false, err
// }
//
// // CheckAuthWithObjectInstances
// func (a Auth) CheckAuthWithObjectInstances(ctx types.ContextParams, targetInsts []inst.Inst) (bool, error) {
//
// 	authAttr := &auth.Attribute{
// 		User: auth.UserInfo{
// 			UserName:   ctx.User,
// 			SupplierID: ctx.SupplierAccount,
// 		},
// 	}
//
// 	for _, objInst := range targetInsts {
// 		instID, err := objInst.GetInstID()
// 		if nil != err {
// 			return false, err
// 		}
// 		authAttr.Resources = append(authAttr.Resources, auth.Resource{
// 			Name:       objInst.GetObject().Object().ObjectName,
// 			InstanceID: uint64(instID),
// 		})
// 	}
//
// 	descision, reason, err := ctx.AuthAPI.Authorize(authAttr)
//
// 	// TODO: need to be implemented, unknown how to use the desision and reason value
// 	_ = descision
// 	_ = reason
//
// 	panic("need to implemented")
//
// 	return true, err
// }
//
// // RegisterAuthResourceWithObject
// func (a Auth) RegisterAuthResourceWithObject(ctx types.ContextParams, targetObjects []model.Object) error {
//
// 	registerMainlineObjectFunc := func(obj model.Object) (*auth.ResourceAttribute, error) {
//
// 		authResourceAttr := &auth.ResourceAttribute{
// 			Object:     obj.Object().ObjectID,
// 			ObjectName: obj.Object().ObjectName,
// 		}
//
// 		// TODO: building the object topology
// 		panic("need to be implemented")
//
// 		return authResourceAttr, nil
// 	}
//
// 	for _, obj := range targetObjects {
//
// 		var authResourceAttr *auth.ResourceAttribute
//
// 		yes, err := obj.IsMainlineObject()
// 		if err != nil {
// 			return err
// 		}
// 		if yes {
// 			authResourceAttr, err = registerMainlineObjectFunc(obj)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			authResourceAttr = &auth.ResourceAttribute{
// 				Object:     obj.Object().ObjectID,
// 				ObjectName: obj.Object().ObjectName,
// 			}
// 		}
//
// 		requestID, err := ctx.AuthAPI.Register(ctx.Context, authResourceAttr)
// 		if nil != err {
// 			return err
// 		}
// 		_ = requestID
//
// 		// TODO: need to be implemented
// 		panic("need to be implemented")
// 	}
//
// 	return nil
// }
//
// // RegisterAuthResourceWithInstance
// func (a Auth) RegisterAuthResourceWithInstance(ctx types.ContextParams, targetInsts []inst.Inst) error {
//
// 	registerMainlineObjectInstFunc := func(objInst inst.Inst) (*auth.ResourceAttribute, error) {
//
// 		// TODO: the data structure of the instance resource needs to be supplemented
// 		authResourceAttr := &auth.ResourceAttribute{}
//
// 		// TODO: building the object topology
// 		panic("need to be implemented")
//
// 		return authResourceAttr, nil
// 	}
//
// 	for _, objInst := range targetInsts {
//
// 		var authResourceAttr *auth.ResourceAttribute
// 		obj := objInst.GetObject()
// 		yes, err := obj.IsMainlineObject()
// 		if err != nil {
// 			return err
// 		}
// 		if yes {
// 			authResourceAttr, err = registerMainlineObjectInstFunc(objInst)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			// TODO: the data structure of the instance resource needs to be supplemented
// 			authResourceAttr = &auth.ResourceAttribute{}
// 		}
//
// 		requestID, err := ctx.AuthAPI.Register(ctx.Context, authResourceAttr)
// 		if nil != err {
// 			return err
// 		}
// 		_ = requestID
//
// 		// TODO: need to be implemented
// 		panic("need to be implemented")
// 	}
//
// 	return nil
// }
//
// // UpdateAuthResourceWithObjecgt
// func (a Auth) UpdateAuthResourceWithObject(ctx types.ContextParams, targetObjects []model.Object) error {
//
// 	registerMainlineObjectFunc := func(obj model.Object) (*auth.ResourceAttribute, error) {
//
// 		// TODO: the data structure of the instance resource needs to be supplemented
// 		authResourceAttr := &auth.ResourceAttribute{}
//
// 		// TODO: building the object topology
// 		panic("need to be implemented")
//
// 		return authResourceAttr, nil
// 	}
//
// 	for _, obj := range targetObjects {
//
// 		var authResourceAttr *auth.ResourceAttribute
// 		yes, err := obj.IsMainlineObject()
// 		if err != nil {
// 			return err
// 		}
// 		if yes {
// 			authResourceAttr, err = registerMainlineObjectFunc(obj)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			// TODO: the data structure of the instance resource needs to be supplemented
// 			authResourceAttr = &auth.ResourceAttribute{}
// 		}
//
// 		requestID, err := ctx.AuthAPI.Update(ctx.Context, authResourceAttr)
// 		if nil != err {
// 			return err
// 		}
// 		_ = requestID
//
// 		// TODO: need to be implemented
// 		panic("need to be implemented")
// 	}
//
// 	return nil
// }
//
// // UpdateAuthResourceWithInstance
// func (a Auth) UpdateAuthResourceWithInstance(ctx types.ContextParams, targetInsts []inst.Inst) error {
//
// 	registerMainlineObjectInstFunc := func(objInst inst.Inst) (*auth.ResourceAttribute, error) {
//
// 		// TODO: the data structure of the instance resource needs to be supplemented
// 		authResourceAttr := &auth.ResourceAttribute{}
//
// 		// TODO: building the object topology
// 		panic("need to be implemented")
//
// 		return authResourceAttr, nil
// 	}
//
// 	for _, objInst := range targetInsts {
//
// 		var authResourceAttr *auth.ResourceAttribute
// 		obj := objInst.GetObject()
// 		yes, err := obj.IsMainlineObject()
// 		if err != nil {
// 			return err
// 		}
// 		if yes {
// 			authResourceAttr, err = registerMainlineObjectInstFunc(objInst)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			// TODO: the data structure of the instance resource needs to be supplemented
// 			authResourceAttr = &auth.ResourceAttribute{}
// 		}
//
// 		requestID, err := ctx.AuthAPI.Update(ctx.Context, authResourceAttr)
// 		if nil != err {
// 			return err
// 		}
// 		_ = requestID
//
// 		// TODO: need to be implemented
// 		panic("need to be implemented")
// 	}
//
// 	return nil
// }
//
// func (a Auth) GetAuthResourceWithObject(ctx types.ContextParams) error {
//
// 	// TODO: the interface is desgining
// 	panic("need to be implemented")
//
// 	return nil
// }
//
// func (a Auth) GetAuthResourceWithInstance(ctx types.ContextParams) error {
//
// 	// TODO: the interface is desgining
// 	panic("beed to be implemented")
//
// 	return nil
// }
//
// // UnRegisterAuthResourceWithObject
// func (a Auth) UnRegisterAuthResourceWithObject(ctx types.ContextParams, targetObjects []model.Object) error {
//
// 	requestID, err := ctx.AuthAPI.Deregister(ctx.Context, &auth.ResourceAttribute{})
// 	_ = requestID
//
// 	// TODO: need to be implemented
// 	panic("need to be implemented")
//
// 	return err
// }
//
// // UnRegisterAuthResourceWithInstance
// func (a Auth) UnRegisterAuthResourceWithInstance(ctx types.ContextParams, targetInsts []inst.Inst) error {
//
// 	requestID, err := ctx.AuthAPI.Deregister(ctx.Context, &auth.ResourceAttribute{})
// 	_ = requestID
//
// 	// TODO: need to be implemented
// 	panic("need to be implemented")
//
// 	return err
// }
