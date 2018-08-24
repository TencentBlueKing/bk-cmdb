import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters('objectModelClassify', {
            $authorizedNavigation: 'authorizedNavigation',
            $classifications: 'classifications'
        }),
        $classify () {
            let $classify = {}
            let relativePath = this.$route.query.relative
            let path = relativePath || this.$route.fullPath
            for (let i = 0; i < this.$authorizedNavigation.length; i++) {
                const classify = this.$authorizedNavigation[i]
                if (classify.hasOwnProperty('path') && classify.path === path) {
                    $classify = classify
                    break
                }
                if (classify.children && classify.children.length) {
                    const targetModel = classify.children.find(model => model.path === path)
                    if (targetModel) {
                        $classify = targetModel
                        break
                    }
                }
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
