# bk-cmdb-ui

> ui of bk-cmdb

## Preparations

- Node.js >= 10.13.0 (LTS), npm >= 6.4.1. see package.json `engines` field.

## Build Setup

``` bash
# install dependencies
npm install

# set dev config -- API_URL in 'builder/config/index.js'
# the API_URL is the address of apiServer and it should start with 'http(s)://', end with '/'
# serve with hot reload at localhost:9090
npm run dev

# build for production with minification
npm run build

# build for production and view the bundle analyzer report
npm run build --report
```
