package repositories

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetAllCategories() []Category {
	// This would normally fetch from a database
	return mock.Categories
}

func GetCategoryByID(id int) []Category {
	for _, cat := range mock.Categories {
		if cat.ID == id {
			return []Category{cat}
		}
	}
	return []Category{}
}

func CreateCategory(name, description string) Category {
	newCategory := Category{
		ID:          len(mock.Categories) + 1,
		Name:        name,
		Description: description,
	}
	mock.Categories = append(mock.Categories, newCategory)
	return newCategory
}

func UpdateCategory(id int, name, description string) *Category {
	for i, cat := range mock.Categories {
		if cat.ID == id {
			mock.Categories[i].Name = name
			mock.Categories[i].Description = description
			return &mock.Categories[i]
		}
	}
	return nil
}

func DeleteCategory(id int) bool{
	for i, cat := range mock.Categories {
		if cat.ID == id {
			mock.Categories = append(mock.Categories[:i], mock.Categories[i+1:]...)
			return true
		}
	}
	return false
}