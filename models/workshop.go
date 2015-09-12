package models

import "time"

type Workshop struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Explanation string    `json:"explanation"`
	Due         time.Time `json:"due"`
	Class       int64     `json:"class"`
}

func (a *Workshop) Delete() error {
	_, err := db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)

	return err
}

func GetAssignment(id int64) (Workshop, error) {
	a := Workshop{
		ID: id,
	}

	var unixTime int64

	row := db.QueryRow("SELECT name, description, explanation, due, class FROM assignments WHERE id = ?", a.ID)
	err := row.Scan(&a.Name, &a.Description, &a.Explanation, &unixTime, &a.Class)
	if err != nil {
		return a, err
	}

	a.Due = time.Unix(0, unixTime)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (a *Workshop) Submissions() ([]Submission, error) {
	var sm []Submission

	rows, err := db.Query("SELECT id, thoughts, team_name FROM submissions WHERE assignment = ?", a.ID)
	if err != nil {
		return sm, err
	}

	for rows.Next() {
		s := Submission{
			ID: a.ID,
		}

		err = rows.Scan(&s.ID, &s.Thoughts, &s.TeamName)
		if err != nil {
			return sm, err
		}

		sm = append(sm, s)
	}

	return sm, nil
}
