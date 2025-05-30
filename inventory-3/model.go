package main

import (
	"database/sql"
	"errors"
)

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

func (m *Movie) getMovie(db *sql.DB) error {
	query := "SELECT name, quantity, price FROM movies WHERE id = ?"
	rows := db.QueryRow(query, m.Id)
	err := rows.Scan(&m.Name, &m.Quantity, &m.Price)
	if err != nil {
		return err
	}
	return err
}

func (m *Movie) createMovie(db *sql.DB) error {
	query := "INSERT INTO movies(name, quantity, price) VALUES(?, ?, ?)"
	result, err := db.Exec(query, m.Name, m.Quantity, m.Price)
	if err != nil {
		return err
	}
	id, err :=  result.LastInsertId()
	if err !=nil {
		return err
	}
	m.Id = int(id)
	return nil
}

func (m *Movie) updateMovie(db *sql.DB) error {
	query := "UPDATE movies SET name = ?, quantity = ?, price = ? WHERE id = ?"
	result, err := db.Exec(query, m.Name, m.Quantity, m.Price, m.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err !=nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("not such rows exists")
	}
	return err
}

func (m *Movie) deleteMovie(db *sql.DB) error {
	query := "DELETE FROM movies WHERE id = ?"
	_, err := db.Exec(query, m.Id)
	return err
}