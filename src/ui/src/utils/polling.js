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
    this.start()
  }

  // 启动轮询
  start() {
    // 避免多次开启同一个轮询器会重复开启多个 timer，所以默认进来就停止之前的 timer。
    this.stop()

    this.isStop = false

    try {
      const pull = () => {
        if (this.isStop) return false
        this.pollingTimer = setTimeout(() => {
          this.callback()
          pull()
        }, this.duration)
      }

      setTimeout(pull, this.duration)
    } catch (e) {
      this.stop()
      throw Error(e)
    }
  }

  // 停止轮询，在没有引用时也会销毁轮询
  stop() {
    this.isStop = true
    this.pollingTimer = null
    clearTimeout(this.pollingTimer)
  }
}

export { Polling }
