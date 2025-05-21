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

package local

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	transactionSessionRedisKeyNamespace = common.BKCacheKeyV3Prefix + "transaction:session:"
	transactionNumberRedisKeyNamespace  = common.BKCacheKeyV3Prefix + "transaction:number:"
	transactionErrorRedisKeyNamespace   = common.BKCacheKeyV3Prefix + "transaction:error:"
)

type sessionKey string

func (s sessionKey) genSessionKey() string {
	return fmt.Sprintf("%s%s", transactionSessionRedisKeyNamespace, string(s))
}

func (s sessionKey) genTxnNumKey() string {
	return fmt.Sprintf("%s%s", transactionNumberRedisKeyNamespace, string(s))
}

func (s sessionKey) genErrKey() string {
	return fmt.Sprintf("%s%s", transactionErrorRedisKeyNamespace, string(s))
}

// TxnErrorType the error type of the transaction, some error type needs to do special operations like retry
type TxnErrorType string

const (
	// UnknownType unknown error type, means the errors that has no specific type, do not have special logic
	UnknownType TxnErrorType = "1"
	// WriteConflictType mongodb write conflict error type, means the transaction conflicts with others, needs to retry
	WriteConflictType TxnErrorType = "2"
)

// ShardingTxnManager is the sharding transaction manager
type ShardingTxnManager struct {
	cache redis.Client
}

// InitTxnManager is to init txn manager, set the redis storage
func (t *ShardingTxnManager) InitTxnManager(r redis.Client) error {
	t.cache = r
	return nil
}

// DB returns the transaction manager for db
func (t *ShardingTxnManager) DB(dbID string) (*TxnManager, error) {
	if dbID == "" {
		return nil, errors.New("db id is not set")
	}

	return &TxnManager{
		dbID:  dbID,
		cache: t.cache,
	}, nil
}

// TxnManager is the transaction manager
type TxnManager struct {
	// dbID is the unique identifier of a db
	dbID  string
	cache redis.Client
}

// GetAllSessionInfos get all session id and txn number infos from redis
func (t *ShardingTxnManager) GetAllSessionInfos(sessionID, initDBID string) ([]SessionInfo, error) {
	ctx := context.Background()

	// get session info for initial session
	txnNumber, err := t.getTxnNumber(sessionID, ctx)
	if err != nil {
		return nil, err
	}
	initSess := make([]SessionInfo, 0)
	if txnNumber > 0 {
		initSess = append(initSess, SessionInfo{TxnNumber: txnNumber, SessionID: sessionID, DBID: initDBID})
	}

	// get session id info from redis
	key := sessionKey(sessionID).genSessionKey()
	sessionIDInfo, err := t.cache.Get(ctx, key).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			return nil, fmt.Errorf("get %s session id info from redis failed, err: %v", sessionID, err)
		}

		// do not have session for other db, get txn number for this session
		return initSess, nil
	}

	sessionIDMap := make(map[string]string)
	if err = json.UnmarshalFromString(sessionIDInfo, &sessionIDMap); err != nil {
		return nil, fmt.Errorf("unmarshal %s session id info %s failed, err: %v", sessionID, sessionIDInfo, err)
	}

	// get all session txn numbers from redis
	sessionInfos := make([]SessionInfo, 0)
	for dbID, dbSessionID := range sessionIDMap {
		txnNumber, err = t.getTxnNumber(dbSessionID, ctx)
		if err != nil {
			return nil, fmt.Errorf("get %s txn number failed, err: %v", dbID, err)
		}
		// skip session id with no txn number
		if txnNumber == 0 {
			continue
		}
		sessionInfos = append(sessionInfos, SessionInfo{TxnNumber: txnNumber, SessionID: dbSessionID, DBID: dbID})
	}

	return append(sessionInfos, initSess...), nil
}

func (t *ShardingTxnManager) getTxnNumber(sessionID string, ctx context.Context) (int64, error) {
	numKey := sessionKey(sessionID).genTxnNumKey()
	txnNum, err := t.cache.Get(ctx, numKey).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("get %s txn number from redis failed, err: %v", sessionID, err)
	}

	txnNumber, err := strconv.ParseInt(txnNum, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s txn number %s failed, err: %v", sessionID, txnNum, err)
	}
	return txnNumber, nil
}

// GenTxnSessionInfo generate the transaction number and session id from redis.
func (t *TxnManager) GenTxnSessionInfo(cap *txnCapableInfo) (*SessionInfo, error) {
	if cap.Timeout == 0 {
		cap.Timeout = common.TransactionDefaultTimeout
	}

	sessionID, err := t.getTxnSessionID(cap)
	if err != nil {
		return nil, err
	}
	key := sessionKey(sessionID).genTxnNumKey()

	pip := t.cache.Pipeline()
	defer pip.Close()

	// we increase by step 1, so that we can calculate how many transaction has already
	// been executed in a same session.
	pip.SetNX(key, 0, cap.Timeout).Result()
	incrBy := pip.IncrBy(key, 1)
	_, err = pip.Exec()
	if err != nil {
		return nil, err
	}
	num := incrBy.Val()

	// value of (num - num/2 +1) is the real transaction number
	// in a distribute session.
	return &SessionInfo{
		TxnNumber: num,
		SessionID: sessionID,
	}, nil
}

// getTxnSessionID get transaction session id for current db
func (t *TxnManager) getTxnSessionID(cap *txnCapableInfo) (string, error) {
	// initial db of this transaction use initial session id
	if cap.isTxnInitialDB {
		return cap.SessionID, nil
	}

	ctx := context.Background()
	key := sessionKey(cap.SessionID).genSessionKey()

	// get session id info from redis
	sessionIDInfo, err := t.cache.Get(ctx, key).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			return "", fmt.Errorf("get %s session id for db %s failed, err: %v", cap.SessionID, t.dbID, err)
		}

		// generate session id and set session id info of current db to redis
		sessionID, err := GenSessionID()
		if err != nil {
			return "", fmt.Errorf("generate %s session id for db %s failed, err: %v", cap.SessionID, t.dbID, err)
		}
		sessionIDInfo = fmt.Sprintf(`{"%s":"%s"}`, t.dbID, sessionID)
		if err = t.cache.Set(ctx, key, sessionIDInfo, cap.Timeout).Err(); err != nil {
			return "", fmt.Errorf("set %s session id info %s failed, err: %v", cap.SessionID, sessionIDInfo, err)
		}

		return sessionID, nil
	}

	sessionIDMap := make(map[string]string)
	if err = json.UnmarshalFromString(sessionIDInfo, &sessionIDMap); err != nil {
		return "", fmt.Errorf("unmarshal %s session id info %s failed, err: %v", cap.SessionID, sessionIDInfo, err)
	}

	sessionID, exists := sessionIDMap[t.dbID]
	if exists {
		return sessionID, nil
	}

	// generate session id and set session id info of current db to redis
	sessionID, err = GenSessionID()
	if err != nil {
		return "", fmt.Errorf("generate %s session id for db %s failed, err: %v", cap.SessionID, t.dbID, err)
	}

	sessionIDMap[t.dbID] = sessionID
	sessionIDInfo, err = json.MarshalToString(sessionIDMap)
	if err != nil {
		return "", fmt.Errorf("marshal %s session id info %v failed, err: %v", cap.SessionID, sessionIDMap, err)
	}

	if err = t.cache.Set(ctx, key, sessionIDInfo, cap.Timeout).Err(); err != nil {
		return "", fmt.Errorf("set %s session id info %s failed, err: %v", cap.SessionID, sessionIDInfo, err)
	}

	return sessionID, nil
}

// RemoveSessionKey remove transaction session key
func (t *ShardingTxnManager) RemoveSessionKey(sessionID string) error {
	key := sessionKey(sessionID).genSessionKey()
	return t.cache.Del(context.Background(), key).Err()
}

// RemoveTxnNumKey remove transaction number key
func (t *ShardingTxnManager) RemoveTxnNumKey(sessionID string) error {
	key := sessionKey(sessionID).genTxnNumKey()
	return t.cache.Del(context.Background(), key).Err()
}

// PrepareCommitOrAbort prepare transaction commit or abort and reload session
func (t *ShardingTxnManager) PrepareCommitOrAbort(cli *mongo.Client, info *SessionInfo) (mongo.Session, error) {
	// create a session client.
	sess, err := cli.StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session failed, err: %v", err)
	}

	// only for changing the transaction status
	err = sess.StartTransaction()
	if err != nil {
		return nil, fmt.Errorf("start transaction failed: %v", err)
	}

	err = CmdbReloadSession(sess, info)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// PrepareTransaction prepare transaction
func (t *TxnManager) PrepareTransaction(cap *txnCapableInfo, cli *mongo.Client) (mongo.Session, error) {
	// create a session client.
	sess, err := cli.StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session failed, err: %v", err)
	}

	// only for changing the transaction status
	err = sess.StartTransaction()
	if err != nil {
		return nil, fmt.Errorf("start transaction %s failed: %v", cap.SessionID, err)
	}

	info, err := t.GenTxnSessionInfo(cap)
	if err != nil {
		return nil, fmt.Errorf("generate txn number failed, err: %v", err)
	}

	// reset the session info with the session id.
	err = CmdbReloadSession(sess, info)
	if err != nil {
		return nil, fmt.Errorf("reload transaction: %s failed, err: %v", cap.SessionID, err)
	}

	return sess, nil
}

// GetTxnContext create a session if the ctx is a session context, and the bool value is true.
// so the caller must check the bool, and use session only when the bool is true.
// otherwise the caller should not use the session, should call the mongodb command directly.
// Note: this function is always used with mongo.CmdbReleaseSession(ctx, sessCtx) to release the session connection.
func (t *TxnManager) GetTxnContext(ctx context.Context, cli *mongo.Client) (context.Context, mongo.Session, bool,
	error) {

	cap, useTxn, err := t.parseTxnInfoFromCtx(ctx)
	if err != nil {
		return ctx, nil, false, err
	}

	if !useTxn {
		// not use transaction, return directly.
		return ctx, nil, false, nil
	}

	session, err := t.PrepareTransaction(cap, cli)
	if err != nil {
		return ctx, nil, true, err
	}

	// prepare the session context, it tells the driver to run this within a transaction.
	sessCtx := CmdbContextWithSession(ctx, session)

	return sessCtx, session, true, nil
}

type txnCapableInfo struct {
	*metadata.TxnCapable
	// isTxnInitialDB defines whether the transaction initial db is the same with current db
	isTxnInitialDB bool
}

// parseTxnInfoFromCtx try to parse transaction info from context,
// it returns the TxnCable, and a bool to indicate whether it's a transaction context or not.
// so the caller can use the returned TxnCapable only when the bool is true. otherwise it will be panic.
func (t *TxnManager) parseTxnInfoFromCtx(txnCtx context.Context) (*txnCapableInfo, bool, error) {
	id := txnCtx.Value(common.TransactionIdHeader)
	if id == nil {
		// do not use transaction, and return directly.
		return nil, false, nil
	}

	txnID, ok := id.(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid transaction id value： %v", id)
	}

	// parse timeout
	ttl := txnCtx.Value(common.TransactionTimeoutHeader)
	if ttl == nil {
		return nil, false, errors.New("transaction timeout value not exist")
	}

	// parse transaction tenant id and check if
	txnTenantID := txnCtx.Value(common.TransactionTenantIDHeader)
	if txnTenantID == nil {
		return nil, false, errors.New("transaction tenant id value not exist")
	}

	txnTenant, exist := tenant.GetTenant(util.GetStrByInterface(txnTenantID))
	if !exist || txnTenant.Status != types.EnabledStatus {
		return nil, false, fmt.Errorf("transaction tenant id %s is invalid", txnTenantID)
	}

	ttlStr, ok := ttl.(string)
	if !ok {
		return nil, false, fmt.Errorf("invalid transaction timeout value: %v", ttl)
	}

	timeout, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil {
		return nil, false, fmt.Errorf("invalid transaction timeout value, parse %v failed, err: %v", ttl, err)
	}

	cap := &txnCapableInfo{
		TxnCapable: &metadata.TxnCapable{
			Timeout:   time.Duration(timeout),
			SessionID: txnID,
		},
		isTxnInitialDB: txnTenant.Database == t.dbID,
	}
	return cap, true, nil
}

// AutoRunWithTxn auto run with transaction
func (t *TxnManager) AutoRunWithTxn(ctx context.Context, cli *mongo.Client, cmd func(ctx context.Context) error) error {
	cap, useTxn, err := t.parseTxnInfoFromCtx(ctx)
	if err != nil {
		return err
	}

	if !useTxn {
		// not use transaction, run command directly.
		return cmd(ctx)
	}

	session, err := t.PrepareTransaction(cap, cli)
	if err != nil {
		return err
	}

	// prepare the session context, it tells the driver to run this within a transaction.
	sessCtx := CmdbContextWithSession(ctx, session)

	// run the command and check error
	err = cmd(sessCtx)
	if err != nil {
		// release the session connection.
		// Attention: do not use session.EndSession() to do this, it will abort the transaction.
		// mongo.CmdbReleaseSession(ctx, session)
		t.setTxnError(sessionKey(cap.SessionID), err)
		return err
	}
	// release the session connection.
	// Attention: do not use session.EndSession() to do this, it will abort the transaction.
	// mongo.CmdbReleaseSession(ctx, session)
	return nil
}

// setTxnError set mongo raw error type to redis, it may be used in scene server to retry this transaction
func (t *TxnManager) setTxnError(sessionID sessionKey, txnErr error) {
	switch {
	case strings.Contains(txnErr.Error(), "WriteConflict"):
		key := sessionID.genErrKey()
		err := t.cache.SetNX(context.Background(), key, string(WriteConflictType), time.Minute*5).Err()
		if err != nil {
			blog.Errorf("set txn error(%v) failed, err: %v, session id: %s", txnErr, err, sessionID)
		}
	default:
	}
}

// GetTxnError get mongo raw error type in redis, the error may be used in scene server to retry this transaction
func (t *ShardingTxnManager) GetTxnError(sessionID string) TxnErrorType {
	key := sessionKey(sessionID).genErrKey()
	errorType, err := t.cache.Get(context.Background(), key).Result()
	if err != nil && !redis.IsNilErr(err) {
		blog.Errorf("get txn error failed, err: %v, session id: %s", err, sessionID)
		return UnknownType
	}

	if len(errorType) == 0 {
		return UnknownType
	}

	return TxnErrorType(errorType)
}

// GenSessionID TODO
func GenSessionID() (string, error) {
	// mongodb driver used this as it's mongodb session id, and we use it too.
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(id[:]), nil
}

// GenTxnCableAndSetHeader TODO
// generate a session id and set it to header.
func GenTxnCableAndSetHeader(header http.Header, opts ...metadata.TxnOption) (*metadata.TxnCapable, error) {
	sessionID, err := GenSessionID()
	if err != nil {
		return nil, fmt.Errorf("generate session id failed, err: %v", err)
	}
	var timeout time.Duration
	if len(opts) != 0 {
		if opts[0].Timeout < 30*time.Second {
			timeout = common.TransactionDefaultTimeout
		} else {
			timeout = opts[0].Timeout
		}
	} else {
		// set default value
		timeout = common.TransactionDefaultTimeout
	}

	header.Set(common.TransactionIdHeader, sessionID)
	header.Set(common.TransactionTimeoutHeader, strconv.FormatInt(int64(timeout), 10))
	header.Set(common.TransactionTenantIDHeader, httpheader.GetTenantID(header))

	cap := metadata.TxnCapable{
		Timeout:   timeout,
		SessionID: sessionID,
	}
	return &cap, nil
}
