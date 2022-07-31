package structures

import (
	"time"

	jwt "github.com/golang-jwt/jwt"
)

// swagger:model Users
type Users struct {
	ID         string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name       string `json:"name" validate:"alpha"`
	Email      string `json:"email"`
	Pass       string `json:"pass" validate:"alphanum"`
	IsVerified bool
	CreatedAt  time.Time
}

// swagger:response Error
type ErrorNotFound struct {
	// Users not found!!
	// in: body
	Message string `json:"error_message"`
}

// swagger:model Organizations
type Organizations struct {
	Org_ID    string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name      string `json:"name"`
	About     string `json:"about"`
	Website   string `json:"website"`
	U_ID      string `json:"user_id" validate:"uuid"`
	CreatedAt time.Time
}

// swagger:model Memberships
type Memberships struct {
	ID        string `json:"pk" gorm:"primaryKey;autoIncrement:false"`
	U_ID      string `json:"user_id" validate:"uuid"`
	Org_ID    string `json:"org_id" validate:"uuid"`
	CreatedAt time.Time
}

// swagger:model Jobs
type Jobs struct {
	ID          string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Org_id      string `json:"org_id" validate:"uuid"`
	Cat_ID      string `json:"cat_id"`
	Designation string `json:"designation"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	CreatedAt   time.Time
}

// swagger:model Category
type Category struct {
	ID   string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Type string `json:"type"`
}

// swagger:model RequiredSkills
type RequiredSkills struct {
	Job_ID string `json:"job_id" validate:"uuid"`
	Skill  string `json:"skill"`
}

// swagger:model Experience
type Experience struct {
	U_ID       string `json:"user_id" validate:"uuid"`
	Experience string `json:"experience"`
}

// swagger:model Skills
type Skills struct {
	U_ID  string `json:"user_id" validate:"uuid"`
	Skill string `json:"skill"`
}

// swagger:model Profile
type Profile struct {
	U_ID       string `json:"user_id" validate:"uuid"`
	Education  string `json:"education"`
	Phone      string `json:"phone" validate:"numeric"`
	Experience []Experience
	Skills     []Skills
}

// swagger:model Applications
type Applications struct {
	Job_ID    string `json:"job_id" validate:"uuid"`
	U_ID      string `json:"user_id" validate:"uuid"`
	Status    string `json:"status"`
	CreatedAt time.Time
}

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Pass  string `json:"pass"`
	jwt.StandardClaims
}
