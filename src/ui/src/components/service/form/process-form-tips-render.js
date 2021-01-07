import ProcessFormTips from './process-form-tips.vue'

export default function (h, { serviceTemplateId, bizId, property }) {
    return h(ProcessFormTips, {
        props: {
            serviceTemplateId,
            bizId,
            property
        }
    })
}
