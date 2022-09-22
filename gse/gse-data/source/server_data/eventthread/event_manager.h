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

#ifndef _GSE_EVENT_MANAGER_H_
#define _GSE_EVENT_MANAGER_H_

#include "datastruct/safe_map.h"
#include <event.h>
#include <event2/event.h>
#include <event2/util.h>
namespace gse {
namespace data {

typedef int64_t EventIndex;
typedef gse::datastruct::SafeMap<EventIndex, struct event*> EventMap;
typedef event_base* EventBase;

typedef void (*EVENT_HANDLER)(evutil_socket_t, short, void*);

class EventManager
{
public:
    EventManager();
    ~EventManager();

    int regEvent(EventBase base, evutil_socket_t fd, EVENT_HANDLER handler, void* arg, bool isEdge = false, bool isPersist = true);
    int regTimer(EventBase base, EVENT_HANDLER handler, void* arg, int second);
    int regMsTimer(EventBase base, EVENT_HANDLER handler, void* arg, int ms);

    int unregEvent(EventIndex index);
    int resetTimer(EventIndex index, int second);
    void clear();
    int activeEvent(EventIndex index);

    static void dropEvent(EventIndex index, struct event* event, void* arg);

private:
    EventIndex m_index;
    EventMap m_eventMap;
};

} // namespace data
} // namespace gse
#endif
