
#include "mongo.h"

bson_t* create_bcon_new_int32(const char *cmd, int32_t val)
{
   return BCON_NEW(cmd, BCON_INT32(val));
}

