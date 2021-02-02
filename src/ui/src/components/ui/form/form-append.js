export default {
    props: ['type', 'property', 'render'],
    render (h) {
        if (typeof this.render === 'function') {
            return this.render(h, {
                type: this.type,
                property: this.property
            })
        }
        return ''
    }
}
