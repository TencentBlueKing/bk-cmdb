export default {
    name: 'render',
    functional: true,
    props: {
        node: Object,
        tpl: Function
    },
    render (h, ct) {
        let titleClass = ct.props.node.selected ? 'node-title node-selected' : 'node-title'
        if (ct.props.tpl) {
            return ct.props.tpl(ct.props.node, ct)
        }
        return (
            <span domPropsInnerHTML={ct.props.node.name} title={ct.props.node.title} class={titleClass}
                style='user-select: none'
                onClick={() => ct.parent.nodeSelected(ct.props.node)}>
            </span>
        )
    }
}
