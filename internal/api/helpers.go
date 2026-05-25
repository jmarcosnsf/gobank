package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmarcosnsf/gobank/internal/services"
)

func (api *Api) getCurrentHolder(ctx context.Context) (uuid.UUID, services.HolderType, error) {
	rawID := api.Sessions.GetString(ctx, "holder_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.UUID{}, "", err
	}
	htype := services.HolderType(api.Sessions.GetString(ctx, "holder_type"))
	return id, htype, nil
}
