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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewAuthCommand())
}

type authConf struct {
	address      string
	appCode      string
	appSecret    string
	resource     string
	resourceFile string
	logv         int32
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
	checkCmd.Flags().StringVar(userName, "user", "", "the name of the user")
	checkCmd.Flags().StringVar(supplierAccount, "supplier-account", "0", "the supplier id that this user belongs to")
	subCmds = append(subCmds, checkCmd)

	for _, subCmd := range subCmds {
		cmd.AddCommand(subCmd)
	}
	conf.addFlags(cmd)

	return cmd
}

func (c *authConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&c.address, "auth-address", "p", "http://iam.service.consul", "auth center addresses, separated by comma")
	cmd.PersistentFlags().StringVarP(&c.appCode, "app-code", "c", "bk_cmdb", "the app code used for authorize")
	cmd.PersistentFlags().StringVarP(&c.appSecret, "app-secret", "s", "", "the app secret used for authorize")
	cmd.PersistentFlags().StringVarP(&c.resource, "resource", "r", "", "the resource for authorize")
	cmd.PersistentFlags().StringVarP(&c.resourceFile, "rsc-file", "f", "", "the resource file path for authorize")
	cmd.PersistentFlags().Int32VarP(&c.logv, "logV", "v", 0, "the log level of request, default request body log level is 4")
}

type authService struct {
	authorize auth.Authorize
	resource  []meta.ResourceAttribute
}

func newAuthService(c *authConf) (*authService, error) {
	blog.SetV(c.logv)
	if c.address == "" {
		return nil, errors.New("auth address must be set")
	}
	if c.appCode == "" {
		return nil, errors.New("app-code must be set")
	}
	if c.appSecret == "" {
		return nil, errors.New("app-secret must be set")
	}
	if c.resource == "" && c.resourceFile == "" {
		return nil, errors.New("resource must be set via resource flag or resource file specified by rsc-file flag")
	}
	addr := strings.Split(c.address, ",")
	for i := range addr {
		if !strings.HasSuffix(addr[i], "/") {
			addr[i] = addr[i] + "/"
		}
	}
	authConf := authcenter.AuthConfig{
		Address:   addr,
		AppCode:   c.appCode,
		AppSecret: c.appSecret,
		SystemID:  authcenter.SystemIDCMDB,
	}
	authorize, err := auth.NewAuthorize(nil, authConf, nil)
	if err != nil {
		return nil, err
	}
	service := &authService{
		authorize: authorize,
	}
	if c.resource != "" {
		err = json.Unmarshal([]byte(c.resource), &service.resource)
		if err != nil {
			return nil, err
		}
	} else {
		resourceFile, err := os.Open(c.resourceFile)
		if err != nil {
			return nil, fmt.Errorf("fail to open file(%s), err(%s)", c.resourceFile, err.Error())
		}
		defer resourceFile.Close()
		resource, err := ioutil.ReadAll(resourceFile)
		if err != nil {
			blog.Errorf("fail to read ata from resource file(%s), err:%s", resourceFile, err.Error())
			return nil, err
		}
		err = json.Unmarshal(resource, &service.resource)
		if err != nil {
			return nil, err
		}
	}
	return service, nil
}

func runAuthRegisterCmd(c *authConf) error {
	srv, err := newAuthService(c)
	if err != nil {
		return err
	}
	err = srv.authorize.RegisterResource(context.Background(), srv.resource...)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stdout, WithBlueColor("Register successful"))
	return nil
}

func runAuthDeregisterCmd(c *authConf) error {
	srv, err := newAuthService(c)
	if err != nil {
		return err
	}
	err = srv.authorize.DeregisterResource(context.Background(), srv.resource...)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stdout, WithBlueColor("Deregister successful"))
	return nil
}

func runAuthUpdateCmd(c *authConf) error {
	srv, err := newAuthService(c)
	if err != nil {
		return err
	}
	for _, res := range srv.resource {
		err = srv.authorize.UpdateResource(context.Background(), &res)
		if err != nil {
			return err
		}

	}
	_, _ = fmt.Fprintln(os.Stdout, WithBlueColor("Update successful"))
	return nil
}

func runAuthCheckCmd(c *authConf, userName string, supplierAccount string) error {
	srv, err := newAuthService(c)
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
		_, _ = fmt.Fprintln(os.Stdout, WithGreenColor("Authorized"))
	} else {
		_, _ = fmt.Fprintln(os.Stdout, WithRedColor("Unauthorized"))
	}
	return nil
}
