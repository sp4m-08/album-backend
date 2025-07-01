package models

type Album struct {
	ID       uint    `json:"id gorm:"primary key"`
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	Price    float64 `json:"price"`
	albumURL string  `json: "imageUrL"`
}

// var Albums = []Album{
// 	{ID: 1, Title: "Blue terrain", Artist: "John Coltrane", Price: 56.78},
// 	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }
