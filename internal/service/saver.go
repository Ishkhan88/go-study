package service

import (
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func StartSaver(ch <-chan model.Entity) {
	for e := range ch {
		repository.SaveEntity(e)
	}
}
