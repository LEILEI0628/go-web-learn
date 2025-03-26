package repository

import (
	"go-web-learn/wire/repository/dao"
)

type TestRepository struct {
	dao *dao.TestDAO
}

func NewTestRepository(dao *dao.TestDAO) *TestRepository {
	return &TestRepository{dao: dao}
}
