package handlers

import (
	"encoding/json"
	"github.com/evgrep/simplyws/internal/data/contract"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"unsafe"
)

type RequestHandler struct {
	dataProvider contract.DataProviderInt
}

func NewRequestHandler(dataProvider contract.DataProviderInt) *RequestHandler {
	requestHandler := &RequestHandler{}
	requestHandler.dataProvider = dataProvider
	return requestHandler
}

func (rh *RequestHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// we support only Get for /entities
		http.NotFound(w, r)
		return
	}

	reqParams := NewGetEntitiesRequestParams()

	// If cannot unmarshall the body into the struct, then return Bad Request
	err := json.NewDecoder(r.Body).Decode(&reqParams)
	if err != nil {
		http.Error(w, "Unable to deserialize the body", http.StatusBadRequest)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(reqParams)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entities, _ := rh.dataProvider.GetEntities(toDBQueryParams(reqParams))

	// write content type json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(entities)

	if err != nil {
		log.Printf("Error occurred while marshalling response into JSON: %v", err)
	}
}

// toDBQueryParams - converts the request parameters struct into the DB query parameters struct
// Since the structures are the same in this case, we just copy one struct into another.
func toDBQueryParams(params *GetEntitiesRequestParams) *contract.GetEntitiesDbParams {
	dbQueryParams := contract.GetEntitiesDbParams{}

	// just copy byte to byte
	p1Slice := (*(*[unsafe.Sizeof(*params)]byte)(unsafe.Pointer(params)))[:]
	p2Slice := (*(*[unsafe.Sizeof(dbQueryParams)]byte)(unsafe.Pointer(&dbQueryParams)))[:]
	copy(p2Slice, p1Slice)

	return &dbQueryParams
}

// GetEntitiesRequestParams - represents the REST API request parameters
type GetEntitiesRequestParams struct {
	IsIncludeLastPrices bool  `json:"include_last_prices"`
	NumberOfLastPrices  int32 `json:"number_of_last_prices" validate:"gte=0"`
	Page                int32 `json:"page" validate:"gte=1"`
	EntitiesPerPage     int32 `json:"entities_per_page" validate:"gte=1"`
}

func NewGetEntitiesRequestParams() *GetEntitiesRequestParams {
	return &GetEntitiesRequestParams{
		EntitiesPerPage:    1,
		NumberOfLastPrices: 3,
	}
}

// ResponseWriter - for the purposes of testing to generate the mock for http/ResponseWriter
//
//go:generate mockery --dir=internal/handlers --name=ResponseWriter --case underscore --output mocks --outpkg=mocks --testonly=true
type ResponseWriter interface {
	http.ResponseWriter
}
