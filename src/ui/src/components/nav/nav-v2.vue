<template>
    <div class="nav-wrapper">
        <v-nav-classify-list
            :classifications="staticNavigationClassify"
            :activeClassify="activeClassifyId">
        </v-nav-classify-list>
        <v-nav-classify-list class="classify-list-custom" v-show="customNavigationClassify.length"
            :classifications="customNavigationClassify"
            :activeClassify="activeClassifyId">
        </v-nav-classify-list>
    </div>
</template>
<script>
    import vNavClassifyList from './nav-classify-list'
    import {mapGetters} from 'vuex'
    export default {
        components: {
            vNavClassifyList
        },
        data () {
            return {
                staticClassify: ['bk_index', 'bk_host_manage', 'bk_organization', 'bk_back_config']
            }
        },
        computed: {
            ...mapGetters('navigation', ['authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'classifyNavigationKey', 'classifyModelSequenceKey']),
            customNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            },
            classifyModelSequence () {
                return this.usercustom[this.classifyModelSequenceKey] || {}
            },
            staticNavigationClassify () {
                const classifies = this.$deepClone(this.authorizedNavigation.filter(classify => this.staticClassify.includes(classify.id)))
                classifies.forEach((classify, index) => {
                    if (this.classifyModelSequence.hasOwnProperty(classify.id) && !['bk_host_manage'].includes(classify.id)) {
                        classify['children'].sort((modelA, modelB) => {
                            return this.getModelSequence(classify, modelA) - this.getModelSequence(classify, modelB)
                        })
                    }
                })
                return classifies
            },
            customNavigationClassify () {
                const classifies = this.$deepClone(this.authorizedNavigation.filter(classify => this.customNavigation.includes(classify.id)))
                classifies.forEach((classify, index) => {
                    if (this.classifyModelSequence.hasOwnProperty(classify.id)) {
                        classify['children'].sort((modelA, modelB) => {
                            return this.getModelSequence(classify, modelA) - this.getModelSequence(classify, modelB)
                        })
                    }
                })
                return classifies
            },
            activeClassify () {
                const path = this.$route.fullPath
                return this.authorizedNavigation.find(classify => classify.children.some(model => model.path === path))
            },
            activeClassifyId () {
                return this.activeClassify ? this.activeClassify.id : this.$route.fullPath === '/index' ? 'bk_index' : null
            },
            activeModel () {
                if (this.activeClassify) {
                    const path = this.$route.fullPath
                    return this.activeClassify.children.find(model => model.path === path)
                }
                return null
            }
        },
        watch: {
            activeModel (activeModel) {
                const index = this.staticNavigationClassify[0] || {}
                let breadcrumbs = [{name: this.$t(index.i18n), path: index.path}]
                if (activeModel) {
                    breadcrumbs.push({
                        name: activeModel.hasOwnProperty('i18n') ? this.$t(activeModel.i18n) : activeModel.name,
                        path: activeModel.path
                    })
                }
                this.$store.commit('main/updateBreadcrumbs', breadcrumbs)
            }
        },
        methods: {
            getModelSequence (classify, model) {
                if (this.classifyModelSequence.hasOwnProperty(classify.id)) {
                    const sequence = this.classifyModelSequence[classify.id]
                    const modelSequence = sequence.indexOf(model.id)
                    return modelSequence === -1 ? classify.children.length : modelSequence
                }
                return classify.children.length
            }
        }
    }
</script>
<style lang="scss" scoped>
    .nav-wrapper{
        padding: 43px 0 0 0;
        width: 110px;
        height: 100%;
        color: #fff;
        background-image: linear-gradient(180deg, #4980d4 0%, #265abb 100%), linear-gradient(#5997eb, #5997eb);
        background-blend-mode: normal, normal;
    }
    .classify-list-custom{
        border-top: 1px solid rgba(228, 231, 234, 0.3);
    }
</style>