const mapping = {
  $eq: 'equal',
  $ne: 'not_equal',
  $in: 'in',
  $nin: 'not_in',
  $lt: 'less',
  $lte: 'less_or_equal',
  $gt: 'greater',
  $gte: 'greater_or_equal',
  $range: 'between',
  $nrange: 'not_between',
  $regex: 'contains'
}

export default operator => mapping[operator]
