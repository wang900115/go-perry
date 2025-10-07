package func06

type UserRepo interface {
	FindByID(id int) (string, error)
}

func GetUsername(repo UserRepo, id int) string {
	name, err := repo.FindByID(id)
	if err != nil {
		return "unknown"
	}
	return name
}
