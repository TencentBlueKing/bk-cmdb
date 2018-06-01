var fs = require('fs')
var os = require('os')
var path = require('path')
var commentMap = require('./commentMap.json')

const options = {
    extensions: /\.(scss|css|js|vue)$/,
    exclude: [],
    separator: os.type() === 'Windows_NT' ? '\\' : '/',
    encoding: 'utf8'
}

function walk (dir) {
    var results = []
    var list = fs.readdirSync(dir)
    list.forEach(function(file) {
        file = dir + options.separator + file
        var stat = fs.statSync(file)
        if (stat && stat.isDirectory()) results = results.concat(walk(file))
        else results.push(file)
    })
    return results
}

function getComment () {
    let comment = fs.readFileSync(path.resolve(__dirname, '../LICENSE.txt'), {encoding: options.encoding})
    if (comment) {
        if (!comment.includes('\n')) {
            comment = `/* ${comment.replace(/\*\//g, "* /")} */\n`
        } else {
            comment = `/*\n * ${comment.replace(/\*\//g, "* /").split("\n").join("\n * ")}\n */\n`
        }
    }
    return comment
}

const comment = getComment()
const commentBuffer = Buffer.from(comment)
if (comment) {
    walk(path.resolve(__dirname, '../src')).forEach(file => {
        if (options.extensions.test(file) && !options.exclude.includes(file)) {
            var relativePath = file.replace(path.resolve(__dirname, '../src'), '').replace(/\//g, '\\')
            if (!commentMap.hasOwnProperty(relativePath)) {
                var contentBuffer = fs.readFileSync(file)
                fs.writeFileSync(file, Buffer.concat([commentBuffer, contentBuffer], commentBuffer.length + contentBuffer.length))
                commentMap[relativePath] = true
            }
        }
    })
    fs.writeFile(path.resolve(__dirname, './commentMap.json'), JSON.stringify(commentMap), (err) => {
        if (err) throw new Error(err)
    })
}