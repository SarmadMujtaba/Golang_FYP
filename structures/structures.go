package structures

import (
	"time"
)

// swagger:model Users
type Users struct {
	// The uuid of a user
	ID string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	// The name of a user
	Name string `json:"name" validate:"alpha"`
	// The email of the user
	Email string `json:"email" validate:"email"`
	// The password of the user
	Pass string `json:"pass" validate:"alphanum"`
	// User created time
	CreatedAt time.Time
}

// swagger:response Error
type ErrorNotFound struct {
	// Users not found!!
	// in: body
	Message string `json:"error_message"`
}

// swagger:model Organizations
type Organizations struct {
	// ID of organization
	Org_ID string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	// Name of organization
	Name string `json:"name"`
	// About organization
	About string `json:"about"`
	// Organization's website
	Website string `json:"website"`
	// Organization's owner's ID
	U_ID string `json:"user_id" validate:"uuid"`
	// Created Time
	CreatedAt time.Time
}

// swagger:model Memberships
type Memberships struct {
	// Membership ID
	ID string `json:"pk" gorm:"primaryKey;autoIncrement:false"`
	// ID of the user to abe added as a member
	U_ID string `json:"user_id" validate:"uuid"`
	// Id of organization against which member is to be added
	Org_ID string `json:"org_id" validate:"uuid"`
	// created time
	CreatedAt time.Time
}

// swagger:model Jobs
type Jobs struct {
	// ID of the job
	ID string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	// ID of the organization posting the job
	Org_id string `json:"org_id" validate:"uuid"`
	// Job Category ID
	Cat_ID string `json:"cat_id"`
	// job Designation or name
	Designation string `json:"designation"`
	// Description of the job
	Description string `json:"description"`
	// Location of the job
	Location string `json:"location"`
	// Estimated salary of the job
	Salary string `json:"salary"`
	// Job creation time
	CreatedAt time.Time
}

// swagger:model Category
type Category struct {
	ID   string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Type string `json:"type"`
}

// swagger:model RequiredSkills
type RequiredSkills struct {
	// Job ID against which skill is to be added.
	Job_ID string `json:"job_id" validate:"uuid"`
	// skill name
	Skill string `json:"skill"`
}

// swagger:model Experience
type Experience struct {
	// User ID
	U_ID string `json:"user_id" validate:"uuid"`
	// Experience Details
	Experience string `json:"experience"`
}

// swagger:model Skills
type Skills struct {
	// User ID
	U_ID string `json:"user_id" validate:"uuid"`
	// Experience Details
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
