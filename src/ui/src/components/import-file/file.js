import { ref } from '@vue/composition-api'
const file = ref(null)
const state = ref(null)
const error = ref(null)
const response = ref(null)

const change = (event) => {
  const { files: [userFile] } = event.target
  file.value = userFile
}

const clear = () => {
  file.value = null
  state.value = null
  error.value = null
  response.value = null
}

const setState = (value) => {
  state.value = value
}

const setError = (value) => {
  error.value = value
}

const setResponse = (value) => {
  response.value = value
}

export default function () {
  return [{ file, state, error, response }, { change, clear, setState, setError, setResponse }]
}
