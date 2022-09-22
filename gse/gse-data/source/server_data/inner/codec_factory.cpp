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
#include "inner/codec_factory.h"
#include "conf/confItem.h"
#include "inner/support_codecs.h"

namespace gse {
namespace data {

Codec* CodecFactory::CreateCodec(int decode_type)
{
    Codec* codec = nullptr;
#ifdef __INNER_CODE__

    switch (decode_type)
    {
    case D_TYPE_TDM_PACKAGE: {
        codec = new TdmPkgCodec();
    }
    break;
    case D_TYPE_GSEDATA_PACKAGE: {
        codec = new GseDataPkgCodec();
    }
    break;
    case D_TYPE_GSEDATA_PACKAGE_V1: {
        codec = new GseDataPkgCodecV1();
    }
    break;
    case D_TYPE_GSEDATA_V1_FOR_TGLOG_PROXY: {
        codec = new GseDataTglogPkgCodecV1();
    }
    break;
    case D_TYPE_ONLY_TRANSPORT: {
        codec = new TransportCodec();
    }
    break;
    case D_TYPE_GSEDATA_PACKAGE_V2: {
        codec = new GeneralProtocalCodec();
    }
    break;

    default:
        LOG_WARN("the type of decode is D_TYPE_UNKNOWN:%d", decode_type);
        break;
    }
#endif
    return codec;
}

CodecFactory::CodecFactory()
{
}

CodecFactory::~CodecFactory()
{
}

} // namespace data
} // namespace gse
