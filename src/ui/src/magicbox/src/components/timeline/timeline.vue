<template>
    <div class="bk-timeline">
        <ul>
            <li class="bk-timeline-dot" :class="makeClass(item)" v-for="(item, index) in list" :key="index">
                <div class="bk-timeline-time" v-if="item.tag !== ''" @click="toggle(item)">{{item.tag}}</div>
                <div class="bk-timeline-content">
                    <template v-if="isNode(item)">
                        <slot :name="'nodeContent' + index">{{nodeContent(item, index)}}</slot>
                    </template>
                    <template v-else>
                        <div v-html="item.content"></div>                    
                    </template>                      
                </div> 
            </li>
        </ul>
    </div>
</template>
<script>
    import { isVNode } from '../../util.js'
    export default {
        name: 'bk-timeline',
        props: {
            list: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                colorReg: /default|primary|warning|success|danger/
            }
        },
        methods: {
            toggle (item) {
                this.$emit('select', item)
            },
            makeClass (item) {
                if (!item.type || !this.colorReg.test(item.type)) {
                    return 'primary'
                }
                return item.type
            },
            isNode (data) {
                if (isVNode(data.content)) {
                    return true
                } else {
                    return false
                }
            },
            nodeContent (data, index) {
                this.$slots[`nodeContent${index}`] = [data.content]
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/timeline.scss';
</style>
