package integration

import (
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
)

type AppDao struct {
	UserDao dao.UserDao
}

type AppRepository struct {
	UserRepository repository.UserRepository
}

type AppService struct {
	UserService dao.UserDao
}

type AppHandler struct {
	UserHandler dao.UserDao
}
