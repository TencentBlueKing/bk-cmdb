export default {
  data() {
    return {
      active: false
    }
  },
  watch: {
    active(active) {
      this.$emit('active-change', active)
    }
  },
  methods: {
    handleToggle(active) {
      this.active = active
    }
  }
}
