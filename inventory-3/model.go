package main

import "database/sql"

type Movie struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getMovies(db *sql.DB) ([]Movie, error) {
	query := "SELECT id, name, quantity, price FROM movies"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	movies := []Movie{}
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.Id, &m.Name, &m.Quantity, &m.Price)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}