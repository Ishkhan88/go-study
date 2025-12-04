package service

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func StartSaver(ctx context.Context, ch <-chan model.Entity) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-ch:
			repository.SaveEntity(e)
		}
	}
}
