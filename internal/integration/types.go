package integration

import (
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"github.com/kgosLj/opsvoid/internal/service"
	"github.com/kgosLj/opsvoid/internal/web/handler"
)

// AppDao 应用程序 DAO
type AppDao struct {
	UserDao dao.UserDao
}

// AppRepository 应用程序 Repository
type AppRepository struct {
	UserRepository repository.UserRepository
}

// AppService 应用程序 Service
type AppService struct {
	UserService service.UserService
}

// AppHandler 应用程序 Handler
type AppHandler struct {
	UserHandler *handler.UserHandler
}
