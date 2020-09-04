export default {
    functional: true,
    name: 'cmdb-form-table-custom-render',
    props: ['row', 'column', '$index', 'render'],
    render (h, { props }) {
        return props.render({ row: props.row, column: props.column, index: props.$index })
    }
}
