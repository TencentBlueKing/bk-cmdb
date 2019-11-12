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

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common/json"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewAuthCommand())
}

type authConf struct {
	address   string
	appCode   string
	appSecret string
	resource  string
}

func NewAuthCommand() *cobra.Command {
	conf := new(authConf)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	subCmds := make([]*cobra.Command, 0)

	subCmds = append(subCmds, &cobra.Command{
		Use:   "register",
		Short: "register resource to auth center",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthRegisterCmd(conf)
		},
	})

	subCmds = append(subCmds, &cobra.Command{
		Use:   "deregister",
		Short: "deregister resource from auth center",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthDeregisterCmd(conf)
		},
	})

	subCmds = append(subCmds, &cobra.Command{
		Use:   "update",
		Short: "update resource in auth center",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthUpdateCmd(conf)
		},
	})

	userName := new(string)
	supplierAccount := new(string)
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "check if user has the authority to operate resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthCheckCmd(conf, *userName, *supplierAccount)
		},
	}
	checkCmd.Flags().StringVar(userName, "user-name", "", "the name of the user")
	checkCmd.Flags().StringVar(supplierAccount, "supplier-account", "0", "the supplier id that this user belongs to")
	subCmds = append(subCmds, checkCmd)

	for _, subCmd := range subCmds {
		cmd.AddCommand(subCmd)
	}
	conf.addFlags(cmd)

	return cmd
}

func (c *authConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.address, "auth-address", "", "auth center addresses, separated by comma")
	cmd.PersistentFlags().StringVar(&c.appCode, "app-code", "", "the app code used for authorize")
	cmd.PersistentFlags().StringVar(&c.appSecret, "app-secret", "", "the app secret used for authorize")
	cmd.PersistentFlags().StringVar(&c.resource, "resource", "", "the resource for authorize")
}

type authService struct {
	authorize auth.Authorize
	resource  []meta.ResourceAttribute
}

func newAuthService(address string, appCode string, appSecret string, resource string) (*authService, error) {
	if address == "" {
		return nil, errors.New("auth-path must be set")
	}
	if appCode == "" {
		return nil, errors.New("app-code must be set")
	}
	if appSecret == "" {
		return nil, errors.New("app-secret must be set")
	}
	if resource == "" {
		return nil, errors.New("resource must be set")
	}
	addr := strings.Split(address, ",")
	for i := range addr {
		if !strings.HasSuffix(addr[i], "/") {
			addr[i] = addr[i] + "/"
		}
	}
	authConf := authcenter.AuthConfig{
		Address:   addr,
		AppCode:   appCode,
		AppSecret: appSecret,
		SystemID:  authcenter.SystemIDCMDB,
	}
	authorize, err := auth.NewAuthorize(nil, authConf, nil)
	if err != nil {
		return nil, err
	}
	service := &authService{
		authorize: authorize,
	}
	err = json.UnmarshalFromString(resource, service.resource)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func runAuthRegisterCmd(c *authConf) error {
	srv, err := newAuthService(c.address, c.appCode, c.appSecret, c.resource)
	if err != nil {
		return err
	}
	err = srv.authorize.RegisterResource(context.Background(), srv.resource...)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stdout, "register auth resource successful")
	return nil
}

func runAuthDeregisterCmd(c *authConf) error {
	srv, err := newAuthService(c.address, c.appCode, c.appSecret, c.resource)
	if err != nil {
		return err
	}
	err = srv.authorize.DeregisterResource(context.Background(), srv.resource...)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stdout, "deregister auth resource successful")
	return nil
}

func runAuthUpdateCmd(c *authConf) error {
	srv, err := newAuthService(c.address, c.appCode, c.appSecret, c.resource)
	if err != nil {
		return err
	}
	for _, res := range srv.resource {
		err = srv.authorize.UpdateResource(context.Background(), &res)
		if err != nil {
			return err
		}

	}
	_, _ = fmt.Fprintln(os.Stdout, "update auth resource successful")
	return nil
}

func runAuthCheckCmd(c *authConf, userName string, supplierAccount string) error {
	srv, err := newAuthService(c.address, c.appCode, c.appSecret, c.resource)
	if err != nil {
		return err
	}
	a := &meta.AuthAttribute{
		Resources: srv.resource,
		User: meta.UserInfo{
			UserName:        userName,
			SupplierAccount: supplierAccount,
		},
	}
	decision, err := srv.authorize.Authorize(context.Background(), a)
	if err != nil {
		return err
	}
	if decision.Authorized {
		_, _ = fmt.Fprintln(os.Stdout, "user has resource permission")
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "user doesn't have resource permission, reason: %s", decision.Reason)
	}
	return nil
}
