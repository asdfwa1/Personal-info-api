package models

type Person struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         uint8  `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type FilterParams struct {
	Name       string
	Surname    string
	Patronymic string
}

type PaginationParams struct {
	Limit  int
	OffSet int
}
