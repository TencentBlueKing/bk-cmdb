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
 

#ifndef _CMDB_MONGO_H_
#define _CMDB_MONGO_H_

#include <mongoc.h>
int64_t bcon_int64(int64_t val);
bson_t* create_bcon_new_int32(const char *cmd, int32_t val);
bool create_collection_index(mongoc_database_t *db,const char* collectionName, bson_t *index, bson_t *reply, bson_error_t *err);
bool get_collection_indexes(mongoc_database_t *db, const char* collectionName, bson_t *reply, bson_error_t *err);

#endif