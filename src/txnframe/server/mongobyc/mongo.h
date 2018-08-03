
#ifndef _CMDB_MONGO_H_
#define _CMDB_MONGO_H_

#include <mongoc.h>
bson_t* create_bcon_new_int32(const char *cmd, int32_t val);

#endif