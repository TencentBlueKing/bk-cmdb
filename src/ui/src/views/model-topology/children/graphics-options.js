export const modelNodeOptions = {
    nodes: {
        shape: 'circle',
        widthConstraint: 62
    },
    interaction: {
        hover: true
    }
}
export const toolNodeOptions = {
    shape: 'box',
    shapeProperties: {
        borderRadius: 12
    },
    heightConstraint: 14,
    scaling: {
        max: 24
    },
    physics: false,
    fixed: {
        x: true,
        y: true
    },
    color: {
        background: 'background: rgba(24, 24, 24, .8)'
    },
    font: {
        color: '#fff'
    }
}
