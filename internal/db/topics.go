package db

type Topic struct {
	ID   string
	Name string
}

func GetTopics() ([]Topic, error) {
	rows, err := doQuery[Topic]("SELECT id, name FROM topics order by id asc")
	if err != nil {
		return nil, err
	}

	return getList[Topic](rows), nil
}
