const hex2grb = (hex) => {
  const rgb = []
  hex = hex.substr(1)
  if (hex.length === 3) {
    hex = hex.replace(/(.)/g, '$1$1')
  }
  hex.replace(/../g, (color) => {
    rgb.push(parseInt(color, 0x10))
  })
  return rgb
}
const getFileExtension = fileName => fileName.substr((~-fileName.lastIndexOf('.') >>> 0) + 2)

const canvas = document.createElement('canvas')

const getBase64Image = (image, color) => {
  const ctx = canvas.getContext('2d')
  canvas.width = image.width
  canvas.height = image.height
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.drawImage(image, 0, 0, image.width, image.height)
  const imageData = ctx.getImageData(0, 0, image.width, image.height)
  const rgbColor = hex2grb(color)
  for (let i = 0; i < imageData.data.length; i += 4) {
    const [r, g, b] = rgbColor
    imageData.data[i] = r
    imageData.data[i + 1] = g
    imageData.data[i + 2] = b
  }
  ctx.putImageData(imageData, 0, 0)
  return canvas.toDataURL(`image/${getFileExtension(image.src)}`)
}

export function svgToImageUrl(image, options) {
  const base64Image = getBase64Image(image, options.iconColor)
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 32 32" x="28" y="25" fill="${options.iconColor}" width="100">
                    <image width="15" xlink:href="${base64Image}"></image>
                </svg>
            </svg>`)}`
}

export function generateObjIcon(image, options) {
  if (image instanceof Image) {
    const base64Image = getBase64Image(image, options.iconColor)
    return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="22" y="5" fill="${options.iconColor}" width="65" >
                        <image width="15" xlink:href="${base64Image}"></image>
                    </svg>
                </svg>`
  }
  options = image
  return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                </svg>`
}
