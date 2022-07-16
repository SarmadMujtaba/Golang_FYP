package structures

import "time"

type Users struct {
	ID        string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name      string `json:"name" validate:"alpha"`
	Email     string `json:"email" validate:"email"`
	Pass      string `json:"pass" validate:"alphanum"`
	CreatedAt time.Time
}

type Organizations struct {
	Org_ID    string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name      string `json:"name"`
	About     string `json:"about"`
	Website   string `json:"website"`
	U_ID      string `json:"user_id" validate:"uuid"`
	CreatedAt time.Time
}

type Memberships struct {
	ID        string `json:"pk" gorm:"primaryKey;autoIncrement:false"`
	U_ID      string `json:"id" validate:"uuid"`
	Org_ID    string `json:"org_id" validate:"uuid"`
	CreatedAt time.Time
	// Users         Users         `gorm:"foreignKey:U_ID;references:ID"`
	// Organizations Organizations `gorm:"foreignKey:Org_ID;references:Org_ID"`
}

type Jobs struct {
	ID          string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Org_id      string `json:"org_id" validate:"uuid"`
	Cat_ID      string `json:"cat_id"`
	Designation string `json:"designation"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
}

type Category struct {
	ID   string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Type string `json:"type"`
}

type RequiredSkills struct {
	ID    string `json:"id"`
	Skill string `json:"skill"`
}
