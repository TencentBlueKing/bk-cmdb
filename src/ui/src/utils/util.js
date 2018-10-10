const HEX_TO_RGB = (hex) => {
    let rgb = []
    hex = hex.substr(1)
    if (hex.length === 3) {
        hex = hex.replace(/(.)/g, '$1$1')
    }
    hex.replace(/../g, function (color) {
        rgb.push(parseInt(color, 0x10))
    })
    return rgb
}
const GET_FILE_EXTENSION = (fileName) => {
    return fileName.substr((~-fileName.lastIndexOf('.') >>> 0) + 2)
}
const GET_BASE_64_IMAGE = (image, color) => {
    let canvas = document.createElement('canvas')
    let ctx = canvas.getContext('2d')
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    canvas.width = image.width
    canvas.height = image.height
    ctx.drawImage(image, 0, 0, image.width, image.height)
    const imageData = ctx.getImageData(0, 0, image.width, image.height)
    const rgbColor = HEX_TO_RGB(color)
    for (let i = 0; i < imageData.data.length; i += 4) {
        imageData.data[i] = rgbColor[0]
        imageData.data[i + 1] = rgbColor[1]
        imageData.data[i + 2] = rgbColor[2]
    }
    ctx.putImageData(imageData, 0, 0)
    return canvas.toDataURL(`image/${GET_FILE_EXTENSION(image.src)}`)
}

export function generateObjIcon (image, options) {
    if (image instanceof Image) {
        const base64Image = GET_BASE_64_IMAGE(image, options.iconColor)
        return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="35" y="-12" fill="${options.iconColor}" width="35" >
                        <image width="15" xlink:href="${base64Image}"></image>
                    </svg>
                    <foreignObject x="0" y="58" width="100%" height="100%">
                        <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                            <div style="color:${options.fontColor};text-align: center;width: 60px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${options.name}</div>
                        </div>
                    </foreignObject>
                </svg>`
    } else {
        options = image
        return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <foreignObject x="0" y="43" width="100%" height="100%">
                        <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                            <div style="color:${options.fontColor};text-align: center;width: 60px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${options.name}</div>
                        </div>
                    </foreignObject>
                </svg>`
    }
}
