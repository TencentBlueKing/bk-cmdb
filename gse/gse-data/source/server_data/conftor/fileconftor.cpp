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

#include "fileconftor.h"
#include "bbx/gse_errno.h"
namespace gse { 
namespace data {

FileConftor::FileConftor()
{
    //
}

FileConftor::~FileConftor()
{
    //
}

int FileConftor::Start()
{
    return GSE_SUCCESS;
}

int FileConftor::Stop()
{
    return GSE_SUCCESS;
}

int FileConftor::CreateConfItem(const std::string& key, std::string& value)
{
    return GSE_SUCCESS;
}

int FileConftor::GetConfItem(const std::string& key, std::string& value, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag)
{
    return GSE_SUCCESS;
}

int FileConftor::SetConfItem(const std::string& key, const std::string& value)
{
    return GSE_SUCCESS;
}
}
}
