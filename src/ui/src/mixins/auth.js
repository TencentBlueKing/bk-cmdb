import { OPERATION } from '@/dictionary/iam-auth'
export default {
    computed: {
        $OPERATION () {
            return Object.freeze(OPERATION)
        }
    }
}
