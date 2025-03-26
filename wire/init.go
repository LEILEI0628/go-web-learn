//go:build wireinject

// 让wire来注入这里的代码
package wire

import (
	"github.com/google/wire"
	"go-web-learn/wire/repository"
	"go-web-learn/wire/repository/dao"
)

func InitTestRepository() *repository.TestRepository {
	// 这个方法传入各个组件的初始化方法
	wire.Build(repository.NewTestRepository, dao.NewTestDAO, InitTestDB)
	return new(repository.TestRepository)
}
