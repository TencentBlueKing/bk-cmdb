import { ref } from '@vue/composition-api'
const current = ref(1)
const next = () => {
  current.value = Math.min(current.value + 1, 2)
}
const previous = () => {
  current.value = Math.max(current.value - 1, 1)
}

const reset = () => {
  current.value = 1
}

export default function () {
  return [current, { next, previous, reset }]
}
