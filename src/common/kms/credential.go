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

// Package kms decrypts the bk-cmdb runtime credential mounted by KMS and maps it onto the existing configuration keys.
package kms

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	kmsdecrypt "github.com/TencentBlueKing/bk-kms-sdk/go/decrypt"

	"configcenter/src/common/blog"
)

const (
	envEnabled        = "BK_KMS_ENABLED"
	envEnvelopeFile   = "BK_KMS_ENVELOPE_FILE"
	envPrivateKeyFile = "BK_KMS_PRIVATE_KEY_FILE"
)

// Credential is the decrypted bk_cmdb_runtime_config document.
// values uses the configuration keys consumed by cmdb. A key is present only
// when the corresponding JSON field exists, including a confirmed empty string.
type Credential struct {
	values    map[string]string
	gseKey    string
	hasGSEKey bool
}

var (
	loadOnce = new(sync.Once)
	cached   *Credential
)

// Enabled reports whether this process should consume the KMS credential.
func Enabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envEnabled))) {
	case "true", "1", "yes", "on":
		return true
	default:
		return false
	}
}

// Load decrypts and parses the runtime credential once per process.
// A nil credential and a nil error means KMS is disabled.
func Load() (*Credential, error) {
	var cachedErr error

	loadOnce.Do(func() {
		if !Enabled() {
			return
		}
		cached, cachedErr = load()
		if cachedErr != nil {
			blog.Errorf("load kms credential failed, err: %v", cachedErr)
			return
		}
		blog.Infof("load kms credential success")
	})

	return cached, cachedErr
}

func load() (*Credential, error) {
	envelopeFile := strings.TrimSpace(os.Getenv(envEnvelopeFile))
	if envelopeFile == "" {
		return nil, fmt.Errorf("%s is empty", envEnvelopeFile)
	}
	privateKeyFile := strings.TrimSpace(os.Getenv(envPrivateKeyFile))
	if privateKeyFile == "" {
		return nil, fmt.Errorf("%s is empty", envPrivateKeyFile)
	}

	envelope, err := readSecret(envelopeFile)
	if err != nil {
		return nil, fmt.Errorf("read kms envelope: %w", err)
	}
	privateKey, err := readSecret(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("read kms private key: %w", err)
	}
	if envelope == "" || privateKey == "" {
		return nil, fmt.Errorf("kms envelope or private key is empty")
	}

	plaintext, err := kmsdecrypt.Decrypt(envelope, privateKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt kms envelope: %w", err)
	}
	return Parse(plaintext)
}

func readSecret(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// runtimeCredential is the bk_cmdb_runtime_config document.
// A nil pointer means the JSON field is absent or null and must not clear configuration.
// A non-nil pointer, including a pointer to "", overlays that configuration key.
type runtimeCredential struct {
	BkappIDSecret *struct {
		Default *struct {
			AppCode   *string `json:"app_code"`
			AppSecret *string `json:"app_secret"`
		} `json:"default"`
	} `json:"bkapp_id_secret"`
	Mongodb *struct {
		Default *userPass `json:"default"`
		Watch   *userPass `json:"watch"`
	} `json:"mongodb"`
	Redis *struct {
		Default  *redisAuth `json:"default"`
		Snapshot *redisAuth `json:"snapshot"`
	} `json:"redis"`
	Zookeeper *struct {
		Default *userPass `json:"default"`
	} `json:"zookeeper"`
	Elasticsearch *struct {
		Default *userPass `json:"default"`
	} `json:"elasticsearch"`
	Kafka *struct {
		Snapshot *userPass `json:"snapshot"`
	} `json:"kafka"`
	Telemetry *struct {
		Monitor *tokenField `json:"monitor"`
		Trace   *tokenField `json:"trace"`
	} `json:"telemetry"`
	Audit *struct {
		Default *tokenField `json:"default"`
	} `json:"audit"`
	Encryption *struct {
		AesGcmKey                       *string `json:"aes_gcm_key"`
		GseClientPrivateKeyBase64       *string `json:"gse_client_private_key_base64"`
		GseAPIServerPrivateKeyPassword  *string `json:"gse_api_server_private_key_password"`
		GseTaskServerPrivateKeyPassword *string `json:"gse_task_server_private_key_password"`
		TracePrivateKeyPassword         *string `json:"trace_private_key_password"`
		ApigatewayPrivateKeyPassword    *string `json:"apigateway_private_key_password"`
		CmdbServerPrivateKeyPassword    *string `json:"cmdb_server_private_key_password"`
	} `json:"encryption"`
}

type userPass struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type redisAuth struct {
	Password         *string `json:"password"`
	SentinelPassword *string `json:"sentinel_password"`
}

type tokenField struct {
	Token *string `json:"token"`
}

// Parse converts the bk_cmdb_runtime_config JSON into configuration keys.
func Parse(data string) (*Credential, error) {
	doc := new(runtimeCredential)
	if err := json.Unmarshal([]byte(data), doc); err != nil {
		return nil, fmt.Errorf("parse kms credential json: %w", err)
	}
	cred := &Credential{values: make(map[string]string)}

	if app := doc.BkappIDSecret; app != nil && app.Default != nil {
		cred.put(app.Default.AppCode, "apiGW.appCode", "auditCenter.appCode")
		cred.put(app.Default.AppSecret, "apiGW.appSecret", "auditCenter.appSecret")
	}
	if db := doc.Mongodb; db != nil {
		cred.putUserPass(db.Default, "mongodb.usr", "mongodb.pwd")
		cred.putUserPass(db.Watch, "watch.usr", "watch.pwd")
	}
	if redis := doc.Redis; redis != nil {
		cred.putRedis(redis.Default, "redis.pwd", "redis.sentinelPwd")
		cred.putRedis(redis.Snapshot, "redis.snap.pwd", "redis.snap.sentinelPwd")
	}
	if zk := doc.Zookeeper; zk != nil {
		cred.putUserPass(zk.Default, "registerServer.usr", "registerServer.pwd")
	}
	if es := doc.Elasticsearch; es != nil {
		cred.putUserPass(es.Default, "es.usr", "es.pwd")
	}
	if kafka := doc.Kafka; kafka != nil {
		cred.putUserPass(kafka.Snapshot, "kafka.snap.user", "kafka.snap.password")
	}
	if telemetry := doc.Telemetry; telemetry != nil {
		if telemetry.Monitor != nil {
			cred.put(telemetry.Monitor.Token, "monitor.accessToken")
		}
		if telemetry.Trace != nil {
			cred.put(telemetry.Trace.Token, "openTelemetry.bkDataToken")
		}
	}
	if audit := doc.Audit; audit != nil && audit.Default != nil {
		cred.put(audit.Default.Token, "auditCenter.token")
	}
	if enc := doc.Encryption; enc != nil {
		cred.put(enc.AesGcmKey, "crypto.aesGcm.key")
		cred.put(enc.GseAPIServerPrivateKeyPassword, "gse.apiServer.password")
		cred.put(enc.GseTaskServerPrivateKeyPassword, "gse.taskServer.password")
		cred.put(enc.TracePrivateKeyPassword, "openTelemetry.tls.password")
		cred.put(enc.ApigatewayPrivateKeyPassword, "apiGW.tls.password")
		cred.put(enc.CmdbServerPrivateKeyPassword, "tls.password")
		if enc.GseClientPrivateKeyBase64 != nil {
			cred.gseKey = *enc.GseClientPrivateKeyBase64
			cred.hasGSEKey = true
		}
	}
	return cred, nil
}

func (c *Credential) put(value *string, keys ...string) {
	if value == nil {
		return
	}
	for _, key := range keys {
		c.values[key] = *value
	}
}

func (c *Credential) putUserPass(auth *userPass, userKey, passKey string) {
	if auth == nil {
		return
	}
	c.put(auth.Username, userKey)
	c.put(auth.Password, passKey)
}

func (c *Credential) putRedis(auth *redisAuth, passKey, sentinelKey string) {
	if auth == nil {
		return
	}
	c.put(auth.Password, passKey)
	c.put(auth.SentinelPassword, sentinelKey)
}
