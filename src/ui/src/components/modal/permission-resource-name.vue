<template>
    <div>
        <div :class="['skeleton', { 'full-width': relations.length > 1 }]" v-if="fetching"></div>
        <div v-else>{{names.join(' / ') || '--'}}</div>
    </div>
</template>

<script>
    import { IAM_VIEWS_INST_NAME } from './permission-resource-name.js'
    export default {
        props: {
            relations: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                fetching: false,
                names: []
            }
        },
        watch: {
            relations (value) {
                this.fetchName()
            }
        },
        created () {
            this.fetchName()
        },
        methods: {
            async fetchName () {
                this.names = []
                const nameReq = []
                this.relations.forEach(relation => {
                    const [type, id] = relation
                    nameReq.push(IAM_VIEWS_INST_NAME[type](this, id, this.relations))
                })
                try {
                    this.fetching = true
                    const result = await Promise.all(nameReq)
                    result.forEach((name, index) => {
                        this.names.push(`【${this.relations[index][2]}】${name}`)
                    })
                } catch (error) {
                    console.error(error)
                } finally {
                    this.fetching = false
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .skeleton {
        width: 50%;
        height: 24px;
        margin: 2px 0;
        background-image: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 37%, #f2f2f2 63%);
        background-size: 400% 100%;
        background-position: 100% 50%;
        animation: skeleton-loading 1.4s ease infinite;

        &.full-width {
            width: 100%;
        }
    }
    @keyframes skeleton-loading {
        0% {
            background-position: 100% 50%;
        }
        100% {
            background-position: 0 50%;
        }
    }
</style>
