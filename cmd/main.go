package main

import (
	"fmt"
	"github.com/evgrep/simplyws/internal/data/providers"
	"github.com/evgrep/simplyws/internal/handlers"
	"github.com/evgrep/simplyws/internal/middleware"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: app <path to sqlite DB>")
	}

	dataProvider, _ := providers.NewSqliteDataProvider(os.Args[1])

	requestHandler := handlers.NewEntitiesRequestHandler(dataProvider)

	apiService := middleware.NewService().
		SetRequestHandler("/entities", requestHandler)

	apiService.Run()
}
