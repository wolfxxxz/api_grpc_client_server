package registry

import (
	"service_user/internal/interface/cache"
	"service_user/internal/interface/controller"
	"service_user/internal/interface/repository"
	"service_user/internal/usecase/interactor"
)

func (r *registry) NewUserController() *controller.UserController {
	userInteractor := interactor.NewUserInteractor(repository.NewUserRepository(r.db, r.log), cache.NewUserCache(r.log, r.redisDB))
	return controller.NewUserController(userInteractor, r.log, r.config)
}
