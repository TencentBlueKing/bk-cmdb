export default {
    computed: {
        $APP () {
            return {
                height: this.$store.state.appHeight
            }
        }
    }
}
