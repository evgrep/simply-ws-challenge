package providers

import (
	"database/sql"
	"encoding/json"
	"github.com/evgrep/simplyws/internal/data/contract"
	"log"
	_ "modernc.org/sqlite"
)

type SqliteDataProvider struct {
	conn *sql.DB
}

func (s *SqliteDataProvider) Connect(path string) error {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	s.conn = conn

	return nil
}

func (s *SqliteDataProvider) GetEntities(params *contract.GetEntitiesDbParams) ([]contract.Entity, error) {
	var query string

	if params.IsIncludeLastPrices {
		query = s.prepareQueryToFetchEntitiesWithLastPrices(query)
	} else {
		query = s.prepareQueryToFetchEntities(query)
	}

	page := params.Page - 1

	rows, err := s.conn.Query(query,
		params.NumberOfLastPrices,
		params.EntitiesPerPage,
		page*params.EntitiesPerPage,
	)

	if err != nil {
		log.Printf("Error occurred while quering the database: %v ", err)
	}

	var data []contract.Entity
	for rows.Next() {
		entity := contract.Entity{}

		var rawLastPrices string
		err = rows.Scan(&entity.Name, &entity.Symbol, &rawLastPrices, &entity.Snowflake)
		if err != nil {
			return nil, err
		}

		// unmarshall the last prices JSON array
		err = json.Unmarshal([]byte(rawLastPrices), &entity.Prices)

		if err != nil {
			log.Printf("Error occurred while unmarshalling JSON data received from the database: %v ", err)
		}

		data = append(data, entity)
	}

	err = rows.Close()

	if err != nil {
		log.Printf("Error occurred while closing the rows: %v ", err)
	}

	return data, nil
}

func (s *SqliteDataProvider) prepareQueryToFetchEntities(query string) string {
	query = `SELECT sc.name, sc.unique_symbol AS symbol, '' AS last_prices, scs.total AS snowflake
				FROM swsCompany sc
				LEFT OUTER JOIN swsCompanyScore scs 
				ON sc.score_id = scs.id
				LIMIT $2
				OFFSET $3`
	return query
}

func (s *SqliteDataProvider) prepareQueryToFetchEntitiesWithLastPrices(query string) string {
	query = `select sc.name, sc.unique_symbol as symbol, IFNULL(lp.last_prices, '') as last_prices, scs.total as snowflake
				from swsCompany sc
				left outer join (
					select company_id, json_group_array(json_object('date', date, 'price', price)) last_prices
					from (
						select row_number() over (PARTITION by company_id) rn, company_id, scpc.date, price 
						from swsCompanyPriceClose scpc
						order by company_id, scpc.date desc
					)
					where rn <= $1
					group by company_id
				) lp
				on sc.id = lp.company_id
				left outer join swsCompanyScore scs 
				ON sc.score_id = scs.id
				limit $2
				OFFSET $3`
	return query
}

func NewSqliteDataProvider(path string) (*SqliteDataProvider, error) {
	dataProvider := &SqliteDataProvider{}
	err := dataProvider.Connect(path)

	if err != nil {
		return nil, err
	}

	return dataProvider, nil
}
