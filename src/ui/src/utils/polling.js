/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

class Polling {
  /**
   * 轮询
   * @param {Function} callback 回调函数
   * @param {number} duration 轮询间隔
   */
  constructor(callback = () => {}, duration = 5000) {
    this.pollingTimer = null // 轮询的 timer
    this.callback = callback
    this.duration = duration
  }

  // 启动轮询
  start() {
    // 阻止重复开启定时器
    if (this.pollingTimer) return false
    try {
      const pull = () => {
        // 停掉以后，避免再往队列里面插入新任务。
        if (!this.pollingTimer) return false
        this.pollingTimer = setTimeout(async () => {
          // 有可能轮询已经停止了，但是还有任务在队列中，所以需要告知正在队列中的任务，你可以停下来了！
          if (!this.pollingTimer) return false
          try {
            await this.callback()
          } catch (err) {
            console.log(err)
          }
          pull()
        }, this.duration)
      }

      this.pollingTimer = setTimeout(pull, this.duration)
    } catch (e) {
      this.stop()
      throw Error(e)
    }
  }

  // 停止轮询
  stop() {
    clearTimeout(this.pollingTimer)
    this.pollingTimer = null
  }
}

export { Polling }
