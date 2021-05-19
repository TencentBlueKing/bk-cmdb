export default {
  data() {
    return {
      active: false
    }
  },
  watch: {
    active(active) {
      this.$emit('active-change', active)
      this.hackEnterEvent()
    }
  },
  methods: {
    handleToggle(active) {
      this.active = active
    },
    hackEnterEvent() {
      if (this.active) {
        window.addEventListener('keyup', this.handleEnter, true)
      } else {
        window.removeEventListener('keyup', this.handleEnter, true)
      }
    },
    handleEnter(event) {
      if (event.key.toLowerCase() !== 'enter') return
      this.$emit('enter', event)
    }
  }
}
