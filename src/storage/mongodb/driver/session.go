/*
* Tencent is pleased to support the open source community by making 蓝鲸 available.
* Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
* Licensed under the MIT License (the ",License"); you may not use this file except
* in compliance with the License. You may obtain a copy of the License at
* http://opensource.org/licenses/MIT
* Unless required by applicable law or agreed to in writing, software distributed under
* the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
* either express or implied. See the License for the specific language governing permissions and
* limitations under the License.
 */

package driver

import (
	"configcenter/src/storage/mongodb"
)

var _ mongodb.Session = (*session)(nil)

type session struct {
	*transaction
	mongocli *client
}

func newSession(mongocli *client) *session {
	return &session{
		mongocli: mongocli,
	}
}

func (s *session) Open() error {

	session, err := s.mongocli.innerClient.StartSession()
	if nil != err {
		return err
	}

	s.transaction = newSessionTransaction(s.mongocli, session)
	return nil
}

func (s *session) Close() error {

	if nil != s.transaction {
		return s.transaction.Close()
	}
	return nil
}

func (s *session) Collection(collName string) mongodb.CollectionInterface {
	return s.transaction.Collection(collName)
}
