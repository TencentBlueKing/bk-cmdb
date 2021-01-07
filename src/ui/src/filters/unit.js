export default function (value, unit) {
    if (!unit || value === '' || value === null || value === undefined) {
        return value
    }
    return String(value) + String(unit)
}
