import { computed } from '@vue/composition-api'

export default function (groups, properties) {
  return computed(() => {
    const result = []
    groups.value.forEach((group) => {
      const matched = properties.value.filter(property => property.bk_property_group === group.bk_group_id)
      result.push({
        group,
        properties: matched
      })
    })
    return result
  })
}
