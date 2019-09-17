import audit from '@/views/audit/router.config'
import business from '@/views/business/router.config'
import businessModel from '@/views/business-model/router.config'
import businessTopology from '@/views/business-topology/router.config'
import customQuery from '@/views/custom-query/router.config'
import eventpush from '@/views/eventpush/router.config'
import history from '@/views/history/router.config'
import hosts from '@/views/hosts/router.config'
import { businessHostDetails, resourceHostDetails } from '@/views/host-details/router.config'
import model from '@/views/model-manage/router.config'
import modelAssociation from '@/views/model-association/router.config'
import modelTopology from '@/views/model-topology/router.config'
import resource from '@/views/resource/router.config'
import generalModel from '@/views/general-model/router.config'
import serviceTemplate from '@/views/service-template/router.config'
import serviceCategory from '@/views/service-category/router.config'
import serviceInstance from '@/views/service-instance/router.config'
import serviceSynchronous from '@/views/business-synchronous/router.config'
import resourceManagement from '@/views/resource-manage/router.config'
import customFields from '@/views/custom-fields/router.config'
import requireBusiness from '@/views/status/require-business'
import setTemplate from '@/views/set-template/router.config'
const flatternViews = views => {
    const flatterned = []
    views.forEach(view => {
        if (Array.isArray(view)) {
            flatterned.push(...view)
        } else {
            flatterned.push(view)
        }
    })
    return flatterned
}

export const businessViews = flatternViews([
    hosts,
    businessHostDetails,
    customQuery,
    businessTopology,
    serviceTemplate,
    serviceCategory,
    serviceInstance,
    serviceSynchronous,
    customFields,
    setTemplate
])

businessViews.forEach(view => {
    view.components = {
        default: view.component,
        requireBusiness: requireBusiness
    }
})

export const resourceViews = flatternViews([
    business,
    resource,
    history,
    resourceHostDetails,
    generalModel,
    eventpush,
    resourceManagement
])

export const modelViews = flatternViews([
    model,
    modelAssociation,
    modelTopology,
    businessModel
])

export const analysisViews = flatternViews([
    audit
])

export default {
    ...businessViews,
    ...resourceViews,
    ...modelViews,
    ...analysisViews
}
