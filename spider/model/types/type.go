package types

type User struct {
	Id        int    `json:"id" gorm:"column:id"`
	StudentID string `json:"student_id" gorm:"column:student_id"`
	Name      string `json:"name" gorm:"column:name"`
	Grade     string `json:"grade" gorm:"column:grade"`
}

