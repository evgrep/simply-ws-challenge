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

	mockDataProvider.On("GetEntities", dbParams).
		Once().
		Return([]contract.Entity{
			{
				Name: "RIO",
			},
			{
				Name: "BHP",
			},
		}, nil)

	requestHandler := RequestHandler{
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
	mockResponseWriter := mocks2.NewMyResponseWriter(t)
	mockResponseWriter.
		On("Header").
		Return(responseHeaders).
		Once()
	mockResponseWriter.
		On("WriteHeader", http.StatusOK)
	mockResponseWriter.
		On("Write", []byte("[{\"name\":\"RIO\",\"symbol\":\"\",\"snowflake\":0},"+
			"{\"name\":\"BHP\",\"symbol\":\"\",\"snowflake\":0}]\n")).
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
				Name: "RIO",
			},
			{
				Name: "BHP",
			},
		}, nil)

	requestHandler := RequestHandler{
		dataProvider: mockDataProvider,
	}

	request := &http.Request{
		Method: "POST",
		Body:   asReader("{}"),
	}

	// create map
	responseHeaders := http.Header{}
	mockResponseWriter := mocks2.NewMyResponseWriter(t)
	mockResponseWriter.
		On("Header").
		Return(responseHeaders).
		Once()
	mockResponseWriter.
		On("WriteHeader", http.StatusOK)
	mockResponseWriter.
		On("Write", []byte("[{\"name\":\"RIO\",\"symbol\":\"\",\"snowflake\":0},{\"name\":\"BHP\",\"symbol\":\"\",\"snowflake\":0}]\n")).
		Return(0, nil)

	// test it now
	requestHandler.Handle(mockResponseWriter, request)

	// assert expectations
	mockDataProvider.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
	assert.Equal(t, responseHeaders.Get("Content-Type"), "application/json")
}
