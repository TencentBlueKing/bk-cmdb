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

export default class CombineRequest {
    static setup (id, callback) {
        return this.getInstance(id, callback)
    }
    constructor (id, callback) {
        this.id = id
        this.timer = null
        this.data = []
        this.callback = callback
        this.promise = null
    }
    add (payload) {
        this.data.push(payload)
        if (!this.timer) {
            this.promise = new Promise(async (resolve, reject) => {
                this.timer = setTimeout(async () => {
                    const result = this.run()
                    resolve(result)
                    this.reset()
                }, 0)
            })
        }
        return this.promise
    }
    run () {
        if (this.callback) {
            return this.callback(this.data)
        }
    }
    reset () {
        clearTimeout(this.timer)
        this.data = []
        this.timer = null
        this.promise = null
    }
    static getInstance (id, callback) {
        const instances = CombineRequest.instances || {}
        if (!instances[id]) {
            instances[id] = new CombineRequest(id, callback)
        }
        CombineRequest.instances = instances
        return instances[id]
    }
}
