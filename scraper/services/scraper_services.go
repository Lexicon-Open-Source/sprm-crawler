package services

import (
	"context"
	"sprm-crawler/common"
	"sprm-crawler/repository"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func UpsertExtraction(ctx context.Context, extractions []repository.Extraction) error {
	context := context.Background()

	tx, err := common.Pool.Begin(context)
	if err != nil {
		return err
	}

	queries := common.Query.WithTx(tx)

	br := queries.UpsertExtraction(ctx, lo.Map(extractions, func(extraction repository.Extraction, _ int) repository.UpsertExtractionParams {
		return repository.UpsertExtractionParams{
			ID:            extraction.ID,
			UrlFrontierID: extraction.UrlFrontierID,
			SiteContent:   extraction.SiteContent,
			ArtifactLink:  extraction.ArtifactLink,
			RawPageLink:   extraction.RawPageLink,
			Language:      extraction.Language,
			PageHash:      extraction.PageHash,
			Metadata:      extraction.Metadata,
			CreatedAt:     extraction.CreatedAt,
			UpdatedAt:     extraction.UpdatedAt,
		}
	}))

	br.Exec(func(int, error) {
		if err != nil {
			log.Error().Err(err).Msg("Error upserting extractions")
		}
	})

	return tx.Commit(ctx)
}
