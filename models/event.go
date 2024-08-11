package models

import (
	"log"
	"time"

	"booking.com/db"
)

type Event struct {
	ID          int64
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID 		int64
}

var events = []Event{}

func (e *Event) Save() error {
	query := `
		INSERT INTO events(name, description, location, dateTime, user_id)
		VALUES (?, ?, ?, ?, ?)`
	
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT id, name, description, location, dateTime, user_id FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Query Error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		var dateTimeStr string
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dateTimeStr, &event.UserID)
		if err != nil {
			log.Printf("Row Scan Error: %v", err)
			return nil, err
		}
	
		event.DateTime, err = time.Parse(time.RFC3339, dateTimeStr) // Use time.RFC3339 for ISO 8601 format
		if err != nil {
			log.Printf("DateTime Parse Error: %v", err)
			return nil, err
		}
	
		events = append(events, event)
	}
	

	return events, nil
}


func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
		UPDATE events
		SET name = ?, description = ?, location = ?, dateTime = ?
		WHERE id = ?	
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}