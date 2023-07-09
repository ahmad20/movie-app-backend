package repositories

import (
	"errors"
	"movie-app-go/entities"
)

type MovieRepository struct {
	data []entities.Movie
}
type MovieRepositoryInterface interface {
	Read(id int) (entities.Movie, error)
	ReadAll() ([]entities.Movie, error)
}

func NewMovieRepository(data []entities.Movie) MovieRepositoryInterface {
	return &MovieRepository{
		data: data,
	}
}

func (repo MovieRepository) Read(id int) (entities.Movie, error) {
	for _, movie := range repo.data {
		if movie.ID == id {
			return movie, nil
		}
	}
	return entities.Movie{}, errors.New("EMPTY_DATA")

}
func (repo MovieRepository) ReadAll() ([]entities.Movie, error) {
	return repo.data, nil
}
