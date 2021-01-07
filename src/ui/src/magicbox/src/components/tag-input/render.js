export default {
    name: 'render',
    functional: true,
    props: {
        node: Object,
        displayKey: String,
        tpl: Function
    },
    render (h, ct) {
        let parentClass = 'bk-selector-node'
        let textClass = 'text'
        if (ct.props.tpl) {
            return ct.props.tpl(ct.props.node, ct)
        }
        return (
            <div class={parentClass}>
                <span domPropsInnerHTML={ct.props.node[ct.props.displayKey]} class={textClass}></span>
            </div>
        )
    }
}
