/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package kms

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ccconfig "configcenter/src/common/core/cc/config"

	"github.com/spf13/viper"

	"configcenter/src/common/blog"
)

// Apply writes credential fields onto configuration keys that already exist.
func (c *Credential) Apply(v *viper.Viper) error {
	if c == nil || v == nil {
		return nil
	}

	overlay := make(map[string]interface{})
	for key, value := range c.values {
		if !v.IsSet(key) {
			continue
		}

		// set the innermost key value
		node := overlay
		path := strings.Split(key, ".")
		for _, k := range path[:len(path)-1] {
			child, ok := node[k].(map[string]interface{})
			if !ok {
				child = make(map[string]interface{})
				node[k] = child
			}
			node = child
		}
		node[path[len(path)-1]] = value
	}

	if len(overlay) > 0 {
		if err := v.MergeConfigMap(overlay); err != nil {
			return fmt.Errorf("merge kms credential: %w", err)
		}
	}

	if !c.hasGSEKey {
		return nil
	}
	return writeGSEKey(v, c.gseKey)
}

// writeGSEKey write gse key info from kms to key file.
func writeGSEKey(v *viper.Viper, keyBase64 string) error {
	if v == nil || !v.IsSet("kms.materialDir") {
		return nil
	}
	dir := strings.TrimSpace(v.GetString("kms.materialDir"))
	if dir == "" {
		return nil
	}
	materialDir := filepath.Clean(dir)
	if materialDir == "" {
		return nil
	}
	info, err := os.Stat(materialDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat kms material dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("kms material dir %s is not a directory", materialDir)
	}

	decoded, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return fmt.Errorf("decode gseCert.cmdbKey: %w", err)
	}
	seen := map[string]struct{}{}
	for _, key := range []string{"gse.apiServer.keyFile", "gse.taskServer.keyFile"} {
		if !v.IsSet(key) {
			continue
		}
		path := v.GetString(key)
		if path == "" || filepath.Dir(filepath.Clean(path)) != materialDir {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		if err = writeFileAtomic(path, decoded); err != nil {
			return fmt.Errorf("write gse key file: %w", err)
		}
		blog.Infof("replaced gseCert.cmdbKey at %s", path)
	}
	return nil
}

func writeFileAtomic(path string, content []byte) (err error) {
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err = tmp.Write(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err = tmp.Chmod(0400); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

// FillZk copies ZooKeeper digest credentials from KMS when those fields exist.
func FillZk(zk *ccconfig.ZkConfig) error {
	if zk == nil {
		return nil
	}

	cred, err := Load()
	if err != nil {
		return err
	}
	if cred == nil {
		return nil
	}

	if user, ok := cred.values["registerServer.usr"]; ok {
		zk.User = user
	}
	if password, ok := cred.values["registerServer.pwd"]; ok {
		zk.Password = password
	}
	return nil
}
