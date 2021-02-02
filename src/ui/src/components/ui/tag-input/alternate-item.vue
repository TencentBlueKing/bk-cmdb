<template>
    <li class="alternate-item"
        :class="{
            highlight: index === tagInput.highlightIndex,
            disabled: disabled && !tagInput.renderList
        }"
        @click.stop
        @mousedown.left.stop="tagInput.handleTagMousedown(tag, disabled)"
        @mouseup.left.stop="tagInput.handleTagMouseup(tag, disabled)">
        <template v-if="tagInput.renderList">
            <render-list
                :tag-input="tagInput"
                :keyword="keyword"
                :tag="tag"
                :disabled="disabled">
            </render-list>
        </template>
        <template v-else>
            <span class="item-name" :title="getTitle()" v-html="getItemContent()"></span>
        </template>
    </li>
</template>
<script>
    import RenderList from './render-list.js'
    export default {
        name: 'alternate-item',
        components: {
            RenderList
        },
        // eslint-disable-next-line
        props: ['tagInput', 'tag', 'keyword', 'index'],
        computed: {
            disabled () {
                return this.tagInput.disabledData.includes(this.tag.value)
            }
        },
        methods: {
            getItemContent () {
                let displayText = this.tag.text || this.tag.value
                if (this.keyword) {
                    displayText = displayText.replace(new RegExp(this.keyword, 'g'), `<span>${this.keyword}</span>`)
                }
                return displayText
            },
            getTitle () {
                return this.tagInput.getDisplayText(this.tag)
            }
        }
    }
</script>
