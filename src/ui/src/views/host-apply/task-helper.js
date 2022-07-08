import md5 from 'md5'

const STORAGE_KEY_PREFIX = 'applytask:'

const genKey = salt => `${STORAGE_KEY_PREFIX}${md5(salt)}`

export const setTask = (id, salt) => {
  const key = genKey(salt)
  localStorage.setItem(key, id)
}

export const getTask = (salt) => {
  const key = genKey(salt)
  return localStorage.getItem(key)
}

export const removeTask = (salt) => {
  const key = genKey(salt)
  localStorage.removeItem(key)
}

export const TASK_STATUS = {
  NEW: 'new',
  WAITING: 'waiting',
  EXECUTING: 'executing',
  FINISHED: 'finished',
  FAIL: 'failure'
}
