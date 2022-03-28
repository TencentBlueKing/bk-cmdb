import { reactive, toRefs } from '@vue/composition-api'
import i18n from '@/i18n'
import useState from './state'
const task = reactive({
  current: null,
  message: null,
  queue: [],
  all: [],
  request: {
    id: Symbol('id')
  },
  iconMapping: {
    error: 'bk-icon icon-close-circle-shape',
    finished: 'bk-icon icon-check-circle-shape',
    pending: 'loading'
  },
  textMapdding: {
    error: i18n.t('失败'),
    finished: i18n.t('已完成'),
    pending: i18n.t('下载中'),
    waiting: i18n.t('等待中')
  }
})
const process = async () => {
  const [state] = useState()
  if (!task.queue.length) return
  try {
    task.message = null
    const [current] = task.queue.splice(0, 1)
    task.current = current
    task.current.state = 'pending'
    await state.submit.value(state, toRefs(task))
    task.current.state = 'finished'
    process()
  } catch (error) {
    task.queue.unshift(task.current)
    task.current.state = 'error'
    task.message = error.message
    console.error(error)
  }
}
const start = () => {
  const [state] = useState()
  const queue = new Array(Math.ceil(state.count.value / state.limit.value)).fill(null)
    .map((_, index) => ({
      name: `${state.bk_obj_id.value}_download_${index + 1}`,
      state: 'waiting',
      page: {
        start: index * state.limit.value,
        limit: state.limit.value
      }
    }))
  task.queue = queue
  task.all = queue.slice()
  process()
}

const reset = () => {
  task.current = null
  task.message = null
  task.queue = []
  task.all = []
}

export default function () {
  return [toRefs(task), { start, reset, process }]
}
