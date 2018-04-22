package inst

type iterator struct {
}

func (cli *iterator) Next() (Inst, error) {
	// TODO:当迭代至最后一条数据后，从数据库中读取下一组实例数据，每组固定条数
	return nil, nil
}
