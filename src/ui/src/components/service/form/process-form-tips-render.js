import ProcessFormTips from './process-form-tips.vue'

export default function (h, { serviceTemplateId }) {
    return h(ProcessFormTips, {
        props: {
            serviceTemplateId
        }
    })
}
