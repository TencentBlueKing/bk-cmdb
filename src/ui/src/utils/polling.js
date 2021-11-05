class Polling {
  constructor(callback = () => {}, duration = 5000) {
    this.pollingTimer = null
    this.callback = callback
    this.duration = duration
    this.start()
  }

  start() {
    try {
      const pull = () => {
        this.pollingTimer = setTimeout(() => {
          this.callback()
          if (!this.pollingTimer) return false
          pull()
        }, this.duration)
      }
      setTimeout(pull, this.duration)
    } catch (e) {
      this.stop()
      throw Error(e)
    }
  }

  stop() {
    this.pollingTimer = null
    clearTimeout(this.pollingTimer)
  }
}

export { Polling }
