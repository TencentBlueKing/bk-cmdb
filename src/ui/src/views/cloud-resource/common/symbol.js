import has from 'has'

class Symbols {
  constructor() {
    this.map = {}
  }

  get(name) {
    if (!has(this.map, name)) {
      this.map[name] = Symbol(name)
    }
    return this.map[name]
  }

  get all() {
    return Object.values(this.map)
  }
}

export default new Symbols()
