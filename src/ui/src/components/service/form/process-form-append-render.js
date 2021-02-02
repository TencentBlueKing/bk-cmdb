import ProcessFormAppend from './process-form-append.vue'

export default function (h, { serviceTemplateId, bizId, property }) {
    return h(ProcessFormAppend, {
        props: {
            serviceTemplateId,
            bizId,
            property
        }
    })
}
