import { setupValidator } from '@/setup/validate'

export default async function (app, to, from) {
    const functions = [setupValidator]
    return Promise.all(functions.map(func => func(app, to, from)))
}
