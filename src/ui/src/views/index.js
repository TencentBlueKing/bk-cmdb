import audit from '@/views/audit/router.config'
import business from '@/views/business/router.config'
import customQuery from '@/views/custom-query/router.config'
import eventpush from '@/views/eventpush/router.config'
import resourceHostDetails from '@/views/host-details/router.config'
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
import operation from '@/views/operation/router.config'

import setSync from '@/views/set-sync/router.config'
import setTemplate from '@/views/set-template/router.config'

import hostApply from '@/views/host-apply/router.config'
import businessTopology from '@/views/business-topology/router.config'

import statusPermission from '@/views/status/permission'
import statusError from '@/views/status/error'

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

const injectStatusComponents = (views, status = ['permission', 'error']) => {
    views.forEach(view => {
        view.components = {
            default: view.component,
            permission: statusPermission,
            error: statusError
        }
    })
    return views
}

export const businessViews = injectStatusComponents(flatternViews([
    customQuery,
    businessTopology,
    serviceTemplate,
    serviceCategory,
    serviceInstance,
    serviceSynchronous,
    customFields,
    setSync,
    setTemplate,
    hostApply
]))

export const resourceViews = injectStatusComponents(flatternViews([
    business,
    resource,
    resourceHostDetails,
    generalModel,
    eventpush,
    resourceManagement
]))

export const modelViews = injectStatusComponents(flatternViews([
    model,
    modelAssociation,
    modelTopology
]))

export const analysisViews = injectStatusComponents(flatternViews([
    audit,
    operation
]))

export default {
    ...businessViews,
    ...resourceViews,
    ...modelViews,
    ...analysisViews
}
