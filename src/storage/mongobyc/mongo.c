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
 

#include "mongo.h"
int64_t bcon_int64(int64_t val)
{
  return BCON_INT64(val);
}
bson_t* create_bcon_new_int32(const char *cmd, int32_t val)
{
   return BCON_NEW(cmd, BCON_INT32(val));
}

bool create_collection_index(mongoc_database_t *db, const char* collectionName, bson_t *index, bson_t *reply, bson_error_t *err)
{
    /* db command format:

      {
        createIndexes: <collection>,
        indexes: [
            {
                key: {
                    <key-value_pair>,
                    <key-value_pair>,
                    ...
                },
                name: <index_name>,
                <option1>,
                <option2>,
                ...
            },
            { ... },
            { ... }
        ],
        writeConcern: { <write concern> }
    }
    */
   bson_t* createIndexes = BCON_NEW ("createIndexes",
                              BCON_UTF8(collectionName),
                              "indexes",
                              "[",
                              BCON_DOCUMENT(index),
                              "]");

   bool ok = mongoc_database_write_command_with_opts(db, createIndexes, NULL /* opts */, reply, err);
   bson_destroy(createIndexes);
   return ok;
}

bool get_collection_indexes(mongoc_database_t *db, const char* collectionName, bson_t *reply, bson_error_t *err)
{
   bson_t* getIndexes = BCON_NEW ("listIndexes",BCON_UTF8(collectionName));
   bool ok = mongoc_database_read_command_with_opts(db, getIndexes, NULL, NULL /* opts */, reply, err);
   bson_destroy(getIndexes);
   return ok;
}


