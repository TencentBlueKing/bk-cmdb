export default function (item, modelId, propertyId) {
    if (modelId === 'host') {
        return item[modelId][propertyId]
    }
    return item[modelId].map(value => value[propertyId]).join(',')
}
