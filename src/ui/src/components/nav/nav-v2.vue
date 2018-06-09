<template>
    <div class="nav-wrapper">
        <v-nav-classify-list
            :classifications="staticNavigationClassify"
            :activeClassify="activeClassify">
        </v-nav-classify-list>
        <v-nav-classify-list class="classify-list-custom" v-show="customNavigationClassify.length"
            :classifications="customNavigationClassify"
            :activeClassify="activeClassify">
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
                const activeClassify = this.authorizedNavigation.find(classify => classify.children.some(model => model.path === path))
                return activeClassify ? activeClassify.id : path === '/index' ? 'bk_index' : null
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