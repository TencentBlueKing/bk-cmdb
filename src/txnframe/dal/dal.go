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

import "github.com/mongodb/mongo-go-driver/mongo"

type DataAccessLayer struct {
	client *mongo.Client
	dbName string
}

func (d *DataAccessLayer) ApplicationBase() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ApplicationBaseCollection),
	}
}

func (d *DataAccessLayer) History() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(HistoryCollection),
	}
}

func (d *DataAccessLayer) HostBase() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(HostBaseCollection),
	}
}

func (d *DataAccessLayer) HostFavourite() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(HostFavouriteCollection),
	}
}

func (d *DataAccessLayer) InstAssociation() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(InstAssociationCollection),
	}
}

func (d *DataAccessLayer) ModuleBase() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ModuleBaseCollection),
	}
}

func (d *DataAccessLayer) ModuleHostConfig() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ModuleHostConfigCollection),
	}
}

func (d *DataAccessLayer) ObjectAssociation() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ObjectAssociationCollection),
	}
}

func (d *DataAccessLayer) ObjectAttribute() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ObjectAttributeCollection),
	}
}

func (d *DataAccessLayer) ObjectClassify() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ObjectClassifyCollection),
	}
}

func (d *DataAccessLayer) ObjectDescription() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ObjectDescriptionCollection),
	}
}

func (d *DataAccessLayer) ObjectBase() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ObjectBaseCollection),
	}
}

func (d *DataAccessLayer) OperationLog() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(OperationLogCollection),
	}
}

func (d *DataAccessLayer) PlatBase() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(PlatBaseCollection),
	}
}

func (d *DataAccessLayer) Privilege() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(PrivilegeCollection),
	}
}

func (d *DataAccessLayer) Proc2Module() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(Proc2ModuleCollection),
	}
}

func (d *DataAccessLayer) Process() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(ProcessCollection),
	}
}

func (d *DataAccessLayer) PropertyGroup() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(PropertyGroupCollection),
	}
}

func (d *DataAccessLayer) BaseSet() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(BaseSetCollection),
	}
}

func (d *DataAccessLayer) Subscription() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(SubscriptionCollection),
	}
}

func (d *DataAccessLayer) System() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(SystemCollection),
	}
}

func (d *DataAccessLayer) Topology() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(TopologyCollection),
	}
}

func (d *DataAccessLayer) UserAPI() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(UserAPICollection),
	}
}

func (d *DataAccessLayer) UserCustom() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(UserCustomCollection),
	}
}

func (d *DataAccessLayer) UserGroup() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(UserGroupCollection),
	}
}

func (d *DataAccessLayer) UserGrpPrivilege() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(UserGroupPrivilegeCollection),
	}
}

func (d *DataAccessLayer) IDGenerator() CollectionInterface {
	return &Collection{
		mgoClient: d.client.Database(d.dbName).Collection(IDGeneratorCollection),
	}
}
