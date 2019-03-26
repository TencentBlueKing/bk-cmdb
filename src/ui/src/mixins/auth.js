export default {
    methods: {
        $isAuthorized (auth) {
            console.log(auth)
            const [ type, action ] = auth.split('.')
            return this.$store.getters['auth/isAuthorized'](type, action)
        }
    }
}
