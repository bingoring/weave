package container

import (
	"weave-module/config"
	"weave-be/internal/application/services"
	"weave-be/internal/domain/repositories"
	domainServices "weave-be/internal/domain/services"
	infraDB "weave-be/internal/infrastructure/database"
	"weave-be/internal/presentation/handlers"
)

// Container holds all application dependencies
type Container struct {
	cfg *config.Config

	// Repositories
	userRepo              repositories.UserRepository
	weaveRepo             repositories.WeaveRepository
	emailVerificationRepo repositories.EmailVerificationRepository

	// Domain Services
	userDomainService domainServices.UserDomainService

	// Application Services (Use Case Based)
	userService *services.UserApplicationService

	// Handlers
	userHandler  *handlers.UserHandler
	oauthHandler *handlers.OAuthHandler
}

// NewContainer creates and initializes the dependency injection container
func NewContainer(cfg *config.Config) *Container {
	container := &Container{
		cfg: cfg,
	}

	container.initializeRepositories()
	container.initializeDomainServices()
	container.initializeApplicationServices()
	container.initializeHandlers()

	return container
}

func (c *Container) initializeRepositories() {
	c.userRepo = infraDB.NewUserRepository()
	c.emailVerificationRepo = infraDB.NewEmailVerificationRepository()
	// TODO: Implement WeaveRepository when ready
	// c.weaveRepo = infraDB.NewWeaveRepository()
}

func (c *Container) initializeDomainServices() {
	c.userDomainService = domainServices.NewUserDomainService(c.userRepo, c.cfg)
}

func (c *Container) initializeApplicationServices() {
	c.userService = services.NewUserApplicationService(c.userRepo, c.weaveRepo, c.userDomainService, c.emailVerificationRepo, c.cfg)
}

func (c *Container) initializeHandlers() {
	c.userHandler = handlers.NewUserHandler(c.userService)
	c.oauthHandler = handlers.NewOAuthHandler(c.userService, c.cfg)
}

// Getters for accessing dependencies
func (c *Container) UserHandler() *handlers.UserHandler {
	return c.userHandler
}

func (c *Container) UserRepository() repositories.UserRepository {
	return c.userRepo
}

func (c *Container) UserDomainService() domainServices.UserDomainService {
	return c.userDomainService
}

func (c *Container) UserService() *services.UserApplicationService {
	return c.userService
}

func (c *Container) OAuthHandler() *handlers.OAuthHandler {
	return c.oauthHandler
}