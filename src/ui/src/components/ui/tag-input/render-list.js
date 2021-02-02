export default {
    name: 'render-list',
    props: ['tagInput', 'tag', 'index', 'keyword', 'disabled'],
    render (h) {
        return this.tagInput.renderList(h, {
            tag: this.tag,
            index: this.index,
            keyword: this.keyword,
            disabled: this.disabled
        })
    }
}
