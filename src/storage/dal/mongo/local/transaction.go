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
    "errors"
    "net/http"
    "time"

    "configcenter/src/common"
    "configcenter/src/common/blog"
    "configcenter/src/storage/dal"
    "configcenter/src/storage/dal/redis"
    "configcenter/src/storage/dal/types"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// NewMgo returns new RDB
func NewTransaction(enableTxn bool, mConf MongoConf, rConf redis.Config) (dal.Transaction, error) {

    if !enableTxn {
        return &Mongo{tm:&TxnManager{enableTransaction: enableTxn},}, nil
    }
    
    connStr, err := connstring.Parse(mConf.URI)
    if nil != err {
        return nil, err
    }

    redisCli, err := redis.NewFromConfig(rConf)
    if err != nil {
        blog.Errorf("new redis client failed, err: %v", err)
        return nil, err
    }

    timeout := time.Duration(mConf.TimeoutSeconds) * time.Second
    conOpt := options.ClientOptions{
        MaxPoolSize:    &mConf.MaxOpenConns,
        MinPoolSize:    &mConf.MaxIdleConns,
        ConnectTimeout: &timeout,
    }

    client, err := mongo.NewClient(options.Client().ApplyURI(mConf.URI), &conOpt)
    if nil != err {
        return nil, err
    }

    if err := client.Connect(context.TODO()); nil != err {
        return nil, err
    }

    return &Mongo{
        dbc:    client,
        dbname: connStr.Database,
        tm:     &TxnManager{
            cache: redisCli,
            enableTransaction: true,
        },
    }, nil
}


// StartTransaction 
// rewrite the ctx and header with transaction related info
func (c *Mongo) StartTransaction(ctx *context.Context, h http.Header, opts ...types.TxnOption) (dal.Transaction, error) {
    if !c.tm.enableTransaction {
        return c, nil
    }
    sess, err := c.dbc.StartSession()
    if err != nil {
        return nil, err
    }

    se := mongo.SessionExposer{}
    info, err := se.GetSessionInfo(c.sess)
    if err != nil {
        return nil, err
    }
    h.Set(common.BKHTTPCCTransactionNumber, info.TxnNumber)
    h.Set(common.BKHTTPCCTxnSessionID, info.SessionID)
    h.Set(common.BKHTTPCCTxnSessionState, info.SessionState)

    context.WithValue(*ctx, common.CCContextKeyJoinOption, types.JoinOption{
        SessionID:    info.SessionID,
        SessionState: info.SessionState,
        TxnNumber:    info.TxnNumber,
    })
    
    nc := Mongo{
        dbc:    c.dbc,
        dbname: c.dbname,
        tm:     c.tm,
        sess: sess,
    }
    
    if len(opts) != 0 {
        if opts[0].Timeout < 5 * time.Second {
            nc.tm.timeout = 5 * time.Second
        } else  {
            nc.tm.timeout = opts[0].Timeout
        }
    } else {
        // set default value
        nc.tm.timeout = 5 * time.Minute
    }

    if err := c.tm.SaveSession(sess); err != nil {
        return nil, err
    }
    
    if err := sess.StartTransaction(); err != nil {
        return nil, err
    }
    return &nc, nil
}


// CommitTransaction 提交事务
func (c *Mongo) CommitTransaction(ctx context.Context) error {
    if !c.tm.enableTransaction {
        return nil
    }
    sess, err := c.chooseSession(ctx)
    if err != nil {
        return err
    }
    if c.hasSession(ctx) {
        defer func() {
            sess.EndSession(ctx)
            if err := c.tm.DeleteSession(sess); err != nil {
                blog.Errorf("delete txn session failed, err: %v", err)
            }
        }()
    }
    return sess.CommitTransaction(ctx)
}

// AbortTransaction 取消事务
func (c *Mongo) AbortTransaction(ctx context.Context) error {
    if !c.tm.enableTransaction {
        return nil
    }
    
    sess, err := c.chooseSession(ctx)
    if err != nil {
        return err
    }
    if c.hasSession(ctx) {
        defer func() {
            sess.EndSession(ctx)
            if err := c.tm.DeleteSession(sess); err != nil {
                blog.Errorf("delete txn session failed, err: %v", err)
            }
        }()
    }
    return sess.AbortTransaction(ctx)
}

// HasSession 判断context里是否有session信息
func (c *Mongo) hasSession(ctx context.Context) bool {
    v, ok := ctx.Value(common.CCContextKeyJoinOption).(types.JoinOption)
    return ok == true && v.SessionID != ""
}

// GetDistributedSession 获取context里用来做分布式事务的session
func (c *Mongo) getDistributedSession(ctx context.Context) (mongo.Session, error) {
    opt, ok := ctx.Value(common.CCContextKeyJoinOption).(types.JoinOption)
    if !ok {
        return nil, errors.New("can't get distributed session, context has no CCContextKeyJoinOption")
    }

    sess, err := c.dbc.StartSession()
    if err != nil {
        return nil, err
    }
    err = sess.StartTransaction()
    if err != nil {
        return nil, err
    }
    err = c.tm.ConvertToSameSession(sess, opt.SessionID)
    if err != nil {
        return nil, err
    }
    return sess, nil
}


// ChooseSession 选择session，优先选择context里用来做分布式事务的session，其次选择自身本地的
func (c *Mongo) chooseSession(ctx context.Context) (mongo.Session, error) {
    var sess mongo.Session
    var err error
    if c.hasSession(ctx) {
        sess, err = c.getDistributedSession(ctx)
        if err != nil {
            return nil, err
        }
    } else if c.sess != nil {
        sess = c.sess
    } else {
        return nil, types.ErrSessionNotStarted
    }
    return sess, nil
}
