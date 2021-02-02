export default {
    name: 'render-tag',
    props: ['tag', 'index'],
    render (h) {
        return this.$parent.renderTag(h, { index: this.index, tag: this.tag })
    }
}
