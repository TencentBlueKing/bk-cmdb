<template>
    <div class="tag-input-alternate-list-wrapper"
        :class="{
            'is-loading': loading
        }"
        :style="wrapperStyle">
        <ul class="alternate-list" ref="alternateList"
            :style="listStyle"
            @scroll="handleScroll">
            <template v-for="(tag, index) in matchedData">
                <template v-if="tag.hasOwnProperty('children')">
                    <!-- eslint-disable vue/require-v-for-key -->
                    <li class="alternate-group"
                        @click.stop
                        @mousedown.left.stop="tagInput.handleGroupMousedown"
                        @mouseup.left.stop="tagInput.handleGroupMouseup">
                        {{`${tag.value || tag.text}(${tag.children.length})`}}
                    </li>
                    <!-- eslint-disable vue/valid-v-for -->
                    <alternate-item v-for="(child, childIndex) in tag.children"
                        ref="alternateItem"
                        :index="getIndex(index, childIndex)"
                        :tag-input="tagInput"
                        :tag="child"
                        :keyword="keyword">
                    </alternate-item>
                </template>
                <!-- eslint-disable vue/valid-v-for -->
                <alternate-item v-else
                    ref="alternateItem"
                    :tag-input="tagInput"
                    :tag="tag"
                    :index="getIndex(index)"
                    :keyword="keyword">
                </alternate-item>
            </template>
        </ul>
        <p class="alternate-empty" v-if="!loading && !matchedData.length">{{tagInput.emptyText}}</p>
    </div>
</template>

<script>
    import AlternateItem from './alternate-item.vue'
    export default {
        components: {
            AlternateItem
        },
        data () {
            return {
                tagInput: null,
                keyword: '',
                next: true,
                loading: true,
                matchedData: []
            }
        },
        computed: {
            wrapperStyle () {
                const style = {}
                if (this.tagInput && this.tagInput.panelWidth) {
                    style.width = parseInt(this.tagInput.panelWidth) + 'px'
                }
                return style
            },
            listStyle () {
                const style = {
                    'max-height': '192px'
                }
                if (this.tagInput) {
                    const maxHeight = parseInt(this.tagInput.listScrollHeight)
                    if (!isNaN(maxHeight)) {
                        style['max-height'] = maxHeight + 'px'
                    }
                }
                return style
            }
        },
        watch: {
            keyword () {
                this.$nextTick(() => {
                    this.$refs.alternateList.scrollTop = 0
                })
            }
        },
        methods: {
            getIndex (index, childIndex = 0) {
                let flattenedIndex = 0
                this.matchedData.slice(0, index).forEach(tag => {
                    if (tag.hasOwnProperty('children')) {
                        flattenedIndex += tag.children.length
                    } else {
                        flattenedIndex += 1
                    }
                })
                return flattenedIndex + childIndex
            },
            handleScroll () {
                if (this.loading || !this.next) return
                const list = this.$refs.alternateList
                const threshold = 32 // 距离底部2条数据
                if ((list.scrollTop + list.clientHeight) > (list.scrollHeight - threshold)) {
                    this.tagInput.search(this.keyword, this.next)
                }
            }
        }
    }
</script>
