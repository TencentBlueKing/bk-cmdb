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

package identifier

/*
  Host identifier event is a mix event of cc_HostBase, cc_ModuleHostConfig and cc_Process events. it's a combination of
these three kinds of events. It has features as follows:
1. it's a virtual event, which represent the three kinds of events mentioned upper.
2. host, host relation and process events is converted to host event and stored to chain node with instance id,
  which is bk_host_id. these event is not really care about the event's order in a event batch operation. and a same
  host identity event(has same host id) are aggregated to only one event in a batch operation.
3. host identifier event do not really store event's detail. the detail is generated when a user(such as gse) call or
  consume a host identifier event. This help us to insert the event's chain nodes more efficiently with this policy.
4. host identifier watch token data structure is not same with the other resources, such as resources in flow.
5. only changes in host fields in needCareHostFields will trigger a host event in cc_HostBase, otherwise,
  this host event will be skip.
6. host identifier's auth resource is redirect to host resource, which means it's authorized with host resource event.
*/
