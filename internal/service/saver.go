package service

import (
	"context"
	"log"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func StartSaver(ctx context.Context, ch <-chan model.Entity) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-ch:
			if !ok {
				return
			}
			if err := repository.SaveEntity(e); err != nil {
				log.Println("save entity error:", err)
			}
		}
	}
}
