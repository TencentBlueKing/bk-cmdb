export default {
    name: 'render-alternate',
    data () {
        return { tagInput: null }
    },
    render (h) {
        return this.tagInput.defaultAlternate(h)
    }
}
