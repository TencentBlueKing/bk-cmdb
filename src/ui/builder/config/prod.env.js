'use strict'
function getCommitId () {
  let commitId = null
  try {
    commitId = require('child_process').execSync('git rev-parse HEAD').toString().trim()
  } catch (error) {
    console.log('Ignore: Get commit id failed')
  }
  return commitId ? JSON.stringify(commitId) : false
}
module.exports = {
  NODE_ENV: '"production"',
  COMMIT_ID: getCommitId()
}
