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

#ifndef _GSE_DATA_JSON_MACRO_H_
#define _GSE_DATA_JSON_MACRO_H_
#include <json/json.h>

namespace gse {
namespace data {

class SafeJson
{
public:
    static std::string GetString(const Json::Value data, const std::string key, const std::string defaultValue)
    {
        std::string value;

        if (!data.isMember(key))
        {
            value = defaultValue;
            return value;
        }

        if (!data[key].isString())
        {
            value = defaultValue;
            return value;
        }

        value = data.get(key, defaultValue).asString();

        return value;
    }

    static int GetInt(const Json::Value data, const std::string key, const int defaultValue)
    {
        int value;

        if (!data.isMember(key))
        {
            value = defaultValue;
            return value;
        }

        if (!data[key].isInt())
        {
            value = defaultValue;
            return value;
        }

        value = data.get(key, defaultValue).asInt();

        return value;
    }

    static bool GetBool(const Json::Value data, const std::string key, const bool defaultValue)
    {
        bool value;

        if (!data.isMember(key))
        {
            value = defaultValue;
            return value;
        }

        if (!data[key].isBool())
        {
            value = defaultValue;
            return value;
        }

        value = data.get(key, defaultValue).asBool();

        return value;
    }

    static double GetDouble(const Json::Value data, const std::string key, const double defaultValue)
    {
        double value;

        if (!data.isMember(key))
        {
            value = defaultValue;
            return value;
        }

        if (!data[key].isDouble())
        {
            value = defaultValue;
            return value;
        }

        value = data.get(key, defaultValue).asDouble();

        return value;
    }

    //    static float GetDouble(const Json::Value data, const std::string key, const float defaultValue)
    //    {
    //        float value;

    //        if (!data.isMember(key))
    //        {
    //            value = defaultValue;
    //            return value;
    //        }

    //        if (!data[key].isFloat())
    //        {
    //            value = defaultValue;
    //            return value;
    //        }

    //        value = data.get(key, defaultValue).asFloat();

    //        return value;
    //    }
};
} // namespace data
} // namespace gse
#endif // JSON_MACRO_H
