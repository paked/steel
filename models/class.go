package models

type Class struct {
	ID          int64
	Name        string
	Description string
	Image       string
}

func NewClass(name, description string) (Class, error) {
	c := Class{
		Name:        name,
		Description: description,
	}

	result, err := db.Exec("INSERT INTO classes (name, description) VALUES (?, ?)", c.Name, c.Description)
	if err != nil {
		return c, err
	}

	c.ID, err = result.LastInsertId()
	if err != nil {
		return c, err
	}

	return c, nil
}
