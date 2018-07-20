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

package dal

import (
	"fmt"

	"configcenter/src/txnframe/client"
	"configcenter/src/txnframe/dal/mongodb"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type DALInterface interface {
	ApplicationBase() mongodb.CollectionInterface
	History() mongodb.CollectionInterface
	HostBase() mongodb.CollectionInterface
	HostFavourite() mongodb.CollectionInterface
	InstAssociation() mongodb.CollectionInterface
	ModuleBase() mongodb.CollectionInterface
	ModuleHostConfig() mongodb.CollectionInterface
	ObjectAssociation() mongodb.CollectionInterface
	ObjectAttribute() mongodb.CollectionInterface
	ObjectClassify() mongodb.CollectionInterface
	ObjectDescription() mongodb.CollectionInterface
	ObjectBase() mongodb.CollectionInterface
	OperationLog() mongodb.CollectionInterface
	PlatBase() mongodb.CollectionInterface
	Privilege() mongodb.CollectionInterface
	Proc2Module() mongodb.CollectionInterface
	Process() mongodb.CollectionInterface
	PropertyGroup() mongodb.CollectionInterface
	BaseSet() mongodb.CollectionInterface
	Subscription() mongodb.CollectionInterface
	System() mongodb.CollectionInterface
	Topology() mongodb.CollectionInterface
	UserAPI() mongodb.CollectionInterface
	UserCustom() mongodb.CollectionInterface
	UserGroup() mongodb.CollectionInterface
	UserGrpPrivilege() mongodb.CollectionInterface
	IDGenerator() mongodb.CollectionInterface
}

func NewMongoDAL(dbName string, mgoCli *mongo.Client, txnClient client.TxnClient) DALInterface {
	return &DataAccessLayer{
		client:    mgoCli,
		dbName:    dbName,
		txnClient: txnClient,
	}
}

type DataAccessLayer struct {
	client    *mongo.Client
	dbName    string
	txnClient client.TxnClient
}

func (d *DataAccessLayer) ApplicationBase() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ApplicationBaseCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ApplicationBaseCollection),
	}
}

func (d *DataAccessLayer) History() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(HistoryCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, HistoryCollection),
	}
}

func (d *DataAccessLayer) HostBase() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(HostBaseCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, HostBaseCollection),
	}
}

func (d *DataAccessLayer) HostFavourite() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(HostFavouriteCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, HostFavouriteCollection),
	}
}

func (d *DataAccessLayer) InstAssociation() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(InstAssociationCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, InstAssociationCollection),
	}
}

func (d *DataAccessLayer) ModuleBase() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ModuleBaseCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ModuleBaseCollection),
	}
}

func (d *DataAccessLayer) ModuleHostConfig() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ModuleHostConfigCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ModuleHostConfigCollection),
	}
}

func (d *DataAccessLayer) ObjectAssociation() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ObjectAssociationCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ObjectAssociationCollection),
	}
}

func (d *DataAccessLayer) ObjectAttribute() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ObjectAttributeCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ObjectAttributeCollection),
	}
}

func (d *DataAccessLayer) ObjectClassify() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ObjectClassifyCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ObjectClassifyCollection),
	}
}

func (d *DataAccessLayer) ObjectDescription() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ObjectDescriptionCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ObjectDescriptionCollection),
	}
}

func (d *DataAccessLayer) ObjectBase() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ObjectBaseCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ObjectBaseCollection),
	}
}

func (d *DataAccessLayer) OperationLog() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(OperationLogCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, OperationLogCollection),
	}
}

func (d *DataAccessLayer) PlatBase() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(PlatBaseCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, PlatBaseCollection),
	}
}

func (d *DataAccessLayer) Privilege() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(PrivilegeCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, PrivilegeCollection),
	}
}

func (d *DataAccessLayer) Proc2Module() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(Proc2ModuleCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, Proc2ModuleCollection),
	}
}

func (d *DataAccessLayer) Process() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(ProcessCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, ProcessCollection),
	}
}

func (d *DataAccessLayer) PropertyGroup() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(PropertyGroupCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, PropertyGroupCollection),
	}
}

func (d *DataAccessLayer) BaseSet() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(BaseSetCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, BaseSetCollection),
	}
}

func (d *DataAccessLayer) Subscription() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(SubscriptionCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, SubscriptionCollection),
	}
}

func (d *DataAccessLayer) System() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(SystemCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, SystemCollection),
	}
}

func (d *DataAccessLayer) Topology() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(TopologyCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, TopologyCollection),
	}
}

func (d *DataAccessLayer) UserAPI() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(UserAPICollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, UserAPICollection),
	}
}

func (d *DataAccessLayer) UserCustom() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(UserCustomCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, UserCustomCollection),
	}
}

func (d *DataAccessLayer) UserGroup() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(UserGroupCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, UserGroupCollection),
	}
}

func (d *DataAccessLayer) UserGrpPrivilege() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(UserGroupPrivilegeCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, UserGroupPrivilegeCollection),
	}
}

func (d *DataAccessLayer) IDGenerator() mongodb.CollectionInterface {
	return &mongodb.Collection{
		MgoClient:   d.client.Database(d.dbName).Collection(IDGeneratorCollection),
		TxnClient:   d.txnClient,
		PreLockPath: fmt.Sprintf("/%s/%s", d.dbName, IDGeneratorCollection),
	}
}
