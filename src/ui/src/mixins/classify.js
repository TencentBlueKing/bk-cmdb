import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters('objectModelClassify', {
            $authorizedNavigation: 'authorizedNavigation',
            $classifications: 'classifications'
        }),
        $classify () {
            let $classify = {}
            let relativePath = this.$route.meta.relative || this.$route.query.relative || null
            let path = relativePath || this.$route.path
            let returnPath = this.$route.meta.returnPath
            let title = this.$route.meta.title
            for (let i = 0; i < this.$authorizedNavigation.length; i++) {
                const classify = this.$authorizedNavigation[i]
                if (classify.hasOwnProperty('path') && classify.path === path) {
                    $classify = classify
                    break
                }
                if (classify.children && classify.children.length) {
                    const targetModel = classify.children.find(child => child.path === path || child.relative === path)
                    if (targetModel) {
                        $classify = targetModel
                        break
                    }
                }
            }
            if (returnPath) {
                Object.assign($classify, { returnPath })
            } else {
                delete $classify.returnPath
            }
            if (title) {
                Object.assign($classify, { title })
            } else {
                delete $classify.title
            }
            return $classify
        },
        $allModels () {
            const allModels = []
            this.$classifications.forEach(classify => {
                classify['bk_objects'].forEach(model => {
                    allModels.push(model)
                })
            })
            return allModels
        },
        $model () {
            let $model = {}
            let $modelId = this.$classify.id
            if ($modelId) {
                const targetModel = this.$allModels.find(model => model['bk_obj_id'] === $modelId)
                if (targetModel) {
                    $model = targetModel
                }
            }
            return $model
        }
    }
}
