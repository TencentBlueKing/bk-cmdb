const fs = require('fs')
const path = require('path')
const iconSelection = require('../src/assets/icon/cc-icon/selection.json')
const iconMap = {}
iconSelection.icons.forEach(icon => {
    iconMap[`icon-${icon.properties.name}`] = icon.properties.code
})
fs.writeFileSync(path.resolve(__dirname, '../src/assets/json/icon-hex-map.json'), JSON.stringify(iconMap), 'utf-8')