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

package iam

var (
    businessParent = Parent{
        SystemID:   systemID,
        ResourceID: Business,
    }
)

// GenerateResourceTypes generate all the resource types registered to IAM.
func GenerateResourceTypes() []ResourceType{
    resourceTypeList := make([]ResourceType, 0)

    // add public resources
    resourceTypeList = append(resourceTypeList, genPublicResources()...)
    
    // add business resources
    resourceTypeList = append(resourceTypeList, genBusinessResources()...)
    
    return resourceTypeList
}

func genBusinessResources() []ResourceType {
    return []ResourceType{
        {
            ID:            BizHostInstance,
            Name:          "业务主机",
            NameEn:        "Business Host",
            Description:   "业务下的机器",
            DescriptionEn: "hosts under a business",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizHostApply,
            Name:          "属性自动应用",
            NameEn:        "Business Host Apply",
            Description:   "自动应用业务主机的属性信息",
            DescriptionEn: "apply business host attribute automatically",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizCustomQuery,
            Name:          "动态分组",
            NameEn:        "Business Custom Query",
            Description:   "根据条件查询主机信息",
            DescriptionEn: "custom query the host instances",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizModel,
            Name:          "业务模型",
            NameEn:        "Business Model",
            Description:   "业务下的模型",
            DescriptionEn: "business's model",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizModelGroup,
            Name:          "业务模型分组",
            NameEn:        "Business Model Group",
            Description:   "对业务下模型进行分类",
            DescriptionEn: "group the business's model",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizProcessServiceInstance,
            Name:          "服务实例",
            NameEn:        "Service Instance",
            Description:   "服务实例",
            DescriptionEn: "service instance",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizProcessServiceCategory,
            Name:          "服务分类",
            NameEn:        "service category",
            Description:   "服务分类",
            DescriptionEn: "service category is to classify service instances",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizSetTemplate,
            Name:          "集群模板",
            NameEn:        "set template",
            Description:   "集群模板",
            DescriptionEn: "set template is used to instantiate a set",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizTopology,
            Name:          "服务拓扑",
            NameEn:        "service topology",
            Description:   "服务拓扑包含了服务拓扑树中所有的相关元素",
            DescriptionEn: "service topology",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        {
            ID:            BizProcessServiceTemplate,
            Name:          "服务模板",
            NameEn:        "service template",
            Description:   "服务模板用于实例化服务实例",
            DescriptionEn: "service template is used to instantiate a service instance ",
            Parents:        []Parent{businessParent},
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
        
    }
}

func genPublicResources() []ResourceType {
    return []ResourceType{
        {
            ID:            Business,
            Name:          "业务",
            NameEn:        "Business",
            Description:   "业务列表",
            DescriptionEn: "all the business in blueking cmdb.",
            Parents:      nil,
            ProviderConfig: ResourceConfig{
                Path: "",
            },
            Version:        0,
        },
    }
}




