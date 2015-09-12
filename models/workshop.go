package models

import "time"

type Workshop struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Explanation string    `json:"explanation"`
	Due         time.Time `json:"due"`
	Class       int64     `json:"class"`
	Prequel     int64
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

func (w *Workshop) CreatePage(title, contents string) (WorkshopPage, error) {
	var p WorkshopPage

	ps, err := w.Pages()
	if err != nil {
		return p, err
	}

	p = WorkshopPage{
		Workshop: w.ID,
		Title:    title,
		Contents: contents,
		Created:  time.Now().UnixNano(),
		Order:    len(ps),
	}

	res, err := db.Exec("INSERT INTO workshop_pages (workshop, contents, title, created, updated, sequence) VALUES (?, ?, ?, ?, ?, ?)", p.Workshop, p.Contents, p.Title, p.Created, p.Updated, p.Order)
	if err != nil {
		return p, err
	}

	p.ID, err = res.LastInsertId()

	return p, err
}

func (w *Workshop) Pages() ([]WorkshopPage, error) {
	var ps []WorkshopPage

	rows, err := db.Query("SELECT id, workshop, contents, title, created, updated, sequence FROM workshop_pages WHERE workshop = ? ORDER BY sequence ASC", w.ID)
	if err != nil {
		return ps, err
	}

	for rows.Next() {
		p := WorkshopPage{}

		err = rows.Scan(&p.ID, &p.Workshop, &p.Contents, &p.Title, &p.Created, &p.Updated, &p.Order)
		if err != nil {
			return ps, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func (w *Workshop) Page(id int64) (WorkshopPage, error) {
	p := WorkshopPage{}

	rows := db.QueryRow("SELECT id, workshop, contents, title, created, updated, sequence FROM workshop_pages WHERE workshop = ? AND id = ? ", w.ID, id)
	err := rows.Scan(&p.ID, &p.Workshop, &p.Contents, &p.Title, &p.Created, &p.Updated, &p.Order)

	return p, err
}

type WorkshopPage struct {
	ID       int64  `json:"id"`
	Workshop int64  `json:"workshop"`
	Contents string `json:"contents"`
	Title    string `json:"title"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
	Order    int    `json:"order"`
}

func (p *WorkshopPage) Edit(title, contents string) error {
	_, err := db.Exec("UPDATE workshop_pages SET contents = ?, title = ? WHERE workshop = ? AND id = ?", contents, title, p.Workshop, p.ID)
	if err != nil {
		return err
	}

	p.Contents = contents

	return nil
}
