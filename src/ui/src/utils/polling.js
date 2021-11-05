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
    this.isStop = true // 轮询停止标识，true 为停止，false 为轮询中。
  }

  // 启动轮询
  start() {
    if (this.pollingTimer) return false
    this.isStop = false

    try {
      const pull = () => {
        // 停掉以后，避免再往队列里面插入新任务。
        if (this.isStop) return false

        this.pollingTimer = setTimeout(() => {
          // 有可能轮询已经停止了，但是还有任务在队列中，所以需要告知正在队列中的任务，你可以停下来了！
          if (this.isStop) return false
          this.callback()
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
    this.pollingTimer = null
    clearTimeout(this.pollingTimer)
    this.isStop = true
  }
}

export { Polling }
