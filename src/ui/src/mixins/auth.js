export default {
    methods: {
        $isAuthorized (auth = '') {
            const [ type, action ] = auth.split('.')
            return this.$store.getters['auth/isAuthorized'](type, action)
        }
    }
}
