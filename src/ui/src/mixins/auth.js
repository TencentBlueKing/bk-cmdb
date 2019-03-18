import { mapGetters } from 'vuex'
import getValue from 'get-value'
export default {
    computed: {
        $auth () {
            return this.$store.getters['auth/auth']
        }
    },
    methods: {
        $hasAuth (type) {
            return getValue(this.$auth, type, { default: false })
        }
    }
}
