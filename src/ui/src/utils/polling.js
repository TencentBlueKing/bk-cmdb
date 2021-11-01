class Polling {
  constructor(callback = () => {}, duration = 5000) {
    this.pollingTimer = null
    this.callback = callback
    this.duration = duration
    this.start()
  }

  start() {
    try {
      if (this.pollingTimer) return false
      this.pollingTimer = setInterval(() => {
        this.callback()
      }, this.duration)
    } catch (e) {
      console.log(e)
      this.stop()
    }
  }

  stop() {
    this.pollingTimer = null
    clearInterval(this.pollingTimer)
  }
}

export { Polling }
