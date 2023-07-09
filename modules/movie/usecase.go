package movie

import (
	"movie-app-go/entities"
	"movie-app-go/repositories"
)

type useCase struct {
	movieRepo repositories.MovieRepositoryInterface
}

type UseCaseInterface interface {
	GetById(id int) (entities.Movie, error)
	GetAll() ([]entities.Movie, error)
}

func NewUseCase(movieRepo repositories.MovieRepositoryInterface) UseCaseInterface {
	return &useCase{
		movieRepo: movieRepo,
	}
}

func (usecase *useCase) GetById(id int) (entities.Movie, error) {
	movie, err := usecase.movieRepo.Read(id)

	if err != nil {
		return entities.Movie{}, err
	}

	return movie, nil
}

func (usecase *useCase) GetAll() ([]entities.Movie, error) {
	admins, _ := usecase.movieRepo.ReadAll()

	return admins, nil
}
