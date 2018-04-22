package common

// Field create a field
type Field interface {
	Eq(val int) Condition
	NotEq(val interface{}) Condition
	Like(val string) Condition
	In(val []interface{}) Condition
	NotIn(val []interface{}) Condition
	Lt(val interface{}) Condition
	Lte(val interface{}) Condition
	Gt(val interface{}) Condition
	Gte(val interface{}) Condition
}

// Field the field object
type field struct {
	condition  Condition
	fieldName  string
	opeartor   string
	fieldValue interface{}
}

// Eqset a filed equal a value
func (cli *field) Eq(val int) Condition {
	cli.opeartor = "$eq"
	cli.fieldValue = val
	return cli.condition
}

// NotEq set a filed equal a value
func (cli *field) NotEq(val interface{}) Condition {
	cli.opeartor = "$ne"
	cli.fieldValue = val
	return cli.condition
}

// Like field like value
func (cli *field) Like(val string) Condition {
	cli.opeartor = "$regex"
	cli.fieldValue = val
	return cli.condition
}

// In in a array
func (cli *field) In(val []interface{}) Condition {
	cli.opeartor = "$in"
	cli.fieldValue = val
	return cli.condition
}

// NotIn not in a array
func (cli *field) NotIn(val []interface{}) Condition {
	cli.opeartor = "$nin"
	cli.fieldValue = val
	return cli.condition
}

// Lt lower than a  value
func (cli *field) Lt(val interface{}) Condition {
	cli.opeartor = "$lt"
	cli.fieldValue = val
	return cli.condition
}

// Lte lower or equal than a value
func (cli *field) Lte(val interface{}) Condition {
	cli.opeartor = "$lte"
	cli.fieldValue = val
	return cli.condition
}

// Gt greater than a value
func (cli *field) Gt(val interface{}) Condition {
	cli.opeartor = "$ge"
	cli.fieldValue = val
	return cli.condition
}

// Gte greater or euqal than a value
func (cli *field) Gte(val interface{}) Condition {
	cli.opeartor = "$gte"
	cli.fieldValue = val
	return cli.condition
}
