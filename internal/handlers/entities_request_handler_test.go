package handlers

import (
	"github.com/evgrep/simplyws/internal/data/contract"
	"github.com/evgrep/simplyws/internal/data/contract/mocks"
	mocks2 "github.com/evgrep/simplyws/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleIncludeLastSet(t *testing.T) {
	// Prepare the test mocks
	mockDataProvider := mocks.NewDataProviderInt(t)

	dbParams := &contract.GetEntitiesDbParams{
		IsIncludeLastPrices: true,
		NumberOfLastPrices:  3,
		EntitiesPerPage:     1,
		Page:                1,
	}

	// simulate data provider GetEntities returning data
	mockDataProvider.On("GetEntities", dbParams).
		Once().
		Return([]contract.Entity{
			{
				Name:      "Rio Tinto",
				Symbol:    "RIO",
				Snowflake: 1,
				Prices: []contract.LastPrice{
					{
						Date:  "2024-06-11",
						Price: "122.91",
					},
					{
						Date:  "2024-06-07",
						Price: "125.31",
					},
				},
			},
			{
				Name:      "BHP Group Ltd",
				Symbol:    "BHP",
				Snowflake: 2,
				Prices: []contract.LastPrice{
					{
						Date:  "2024-06-11",
						Price: "43.74",
					},
					{
						Date:  "2024-06-07",
						Price: "44.53",
					},
				},
			},
		}, nil)

	requestHandler := EntitiesRequestHandler{
		dataProvider: mockDataProvider,
	}

	request := &http.Request{
		Method: "POST",
		Body: asReader(`{
					"include_last_prices": true
				}`),
	}

	// create map
	responseHeaders := http.Header{}
	mockResponseWriter := mocks2.NewResponseWriter(t)
	mockResponseWriter.
		On("Header").
		Return(responseHeaders).
		Once()
	mockResponseWriter.
		On("WriteHeader", http.StatusOK)
	mockResponseWriter.
		On("Write", []byte("[{\"name\":\"Rio Tinto\",\"symbol\":\"RIO\",\"snowflake\":1,\"last_prices\":"+
			"[{\"date\":\"2024-06-11\",\"price\":122.91},{\"date\":\"2024-06-07\",\"price\":125.31}]},{\"name\":\"BH"+
			"P Group Ltd\",\"symbol\":\"BHP\",\"snowflake\":2,\"last_prices\":[{\"date\":\"2024-06-11\",\"price\":43"+
			".74},{\"date\":\"2024-06-07\",\"price\":44.53}]}]\n")).
		Return(0, nil)

	// test it now
	requestHandler.Handle(mockResponseWriter, request)

	// assert expectations
	mockDataProvider.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
	assert.Equal(t, responseHeaders.Get("Content-Type"), "application/json")
}

func asReader(contents string) io.ReadCloser {
	stringReader := strings.NewReader(contents)
	body := io.NopCloser(stringReader)
	return body
}

func TestHandleIncludeLastUnset(t *testing.T) {
	// Prepare the test mocks
	mockDataProvider := mocks.NewDataProviderInt(t)

	dbParams := &contract.GetEntitiesDbParams{
		IsIncludeLastPrices: false,
		NumberOfLastPrices:  3,
		EntitiesPerPage:     1,
		Page:                1,
	}

	mockDataProvider.On("GetEntities", dbParams).
		Once().
		Return([]contract.Entity{
			{
				Name:      "Rio Tinto",
				Symbol:    "RIO",
				Snowflake: 1,
			},
			{
				Name:      "BHP Group Ltd",
				Symbol:    "BHP",
				Snowflake: 2,
			},
		}, nil)

	requestHandler := EntitiesRequestHandler{
		dataProvider: mockDataProvider,
	}

	request := &http.Request{
		Method: "POST",
		Body:   asReader("{}"),
	}

	// create map
	responseHeaders := http.Header{}
	mockResponseWriter := mocks2.NewResponseWriter(t)
	mockResponseWriter.
		On("Header").
		Return(responseHeaders).
		Once()
	mockResponseWriter.
		On("WriteHeader", http.StatusOK)
	mockResponseWriter.
		On("Write", []byte("[{\"name\":\"Rio Tinto\",\"symbol\":\"RIO\",\"snowflake\":1},{\"name\":"+
			"\"BHP Group Ltd\",\"symbol\":\"BHP\",\"snowflake\":2}]\n")).
		Return(0, nil)

	// test it now
	requestHandler.Handle(mockResponseWriter, request)

	// assert expectations
	mockDataProvider.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
	assert.Equal(t, responseHeaders.Get("Content-Type"), "application/json")
}
