package contract

//go:generate mockery --name=DataProviderInt --case underscore --output mocks --outpkg=mocks
type DataProviderInt interface {
	GetEntities(params *GetEntitiesDbParams) ([]Entity, error)
}

// GetEntitiesDbParams - represents the parameters passed to GetEntities method.
type GetEntitiesDbParams struct {
	IsIncludeLastPrices bool  `json:"include_last_prices"`
	NumberOfLastPrices  int32 `json:"number_of_last_prices" validate:"gte=0"`
	Page                int32 `json:"page" validate:"gte=0"`
	EntitiesPerPage     int32 `json:"entities_per_page" validate:"gte=1"`
}
