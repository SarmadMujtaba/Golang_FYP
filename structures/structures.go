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

// swagger:model UserWrapper
type UserWrapper struct {
	// The name of a user
	Name string `json:"name" validate:"alpha"`
	// The email of the user
	Email string `json:"email" validate:"email"`
	// The password of the user
	Pass string `json:"pass" validate:"alphanum"`
}

// swagger:parameters post-user
type AddUserSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body UserWrapper `json:"body"`
}

// swagger:response Error
type ErrorNotFound struct {
	// Users not found!!
	// in: headers
}

// swagger:parameters userParam deleteParam deleteJob
type Param struct {
	// User ID
	// in: query
	ID string `json:"id"`
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

// swagger:model OrgWrapper
type OrgWrapper struct {
	// The name of a user
	Name string `json:"name" validate:"alpha"`
	// The email of the user
	About string `json:"about"`
	// The password of the user
	Website string `json:"pass" validate:"alphanum"`
	// Organization's owner's ID
	U_ID string `json:"user_id" validate:"uuid"`
}

// swagger:parameters post-org
type AddOrgSwagger struct {
	//  Add details of the Organization
	//  in: body
	//  required: true
	Body OrgWrapper `json:"body"`
}

// swagger:parameters orgParam deleteOrgParam
type OrgParam struct {
	// Organization ID
	// in: query
	ID string `json:"id" validate:"uuid"`
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

// swagger:model memberWrapper
type MemberWrapper struct {
	// User ID
	// in: body
	U_ID string `json:"user_id" validate:"uuid"`
	// Organization ID
	// in: body
	Org_ID string `json:"org_id" validate:"uuid"`
}

// swagger:parameters post-member
type AddMemberSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body MemberWrapper `json:"body"`
}

// swagger:parameters memberParam
type MemberWrapper2 struct {
	// Organization ID
	// in: query
	Org_ID string `json:"org_id" validate:"uuid"`
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

// swagger:model jobWrapper
type JobWrapper struct {
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
}

// swagger:parameters post-job
type AddJobSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body JobWrapper `json:"body"`
}

//swagger:parameters jobParam
type JobParam struct {
	// Organization ID
	// in: query
	ID string `json:"org_id"`
}

//swagger:parameters nameParam
type JobNameParam struct {
	// Organization ID
	// in: query
	Job_Name string `json:"job_name"`
}

type Category struct {
	ID   string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Type string `json:"type"`
}

//swagger:parameters catParam
type CatParam struct {
	// Category ID
	// in: query
	Category_ID string `json:"category_id"`
}

// swagger:model RequiredSkills
type RequiredSkills struct {
	// Job ID against which skill is to be added.
	Job_ID string `json:"job_id" validate:"uuid"`
	// skill name
	Skill string `json:"skill"`
}

// swagger:parameters post-requiredSkill
type AddRequiredSkillsSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body RequiredSkills `json:"body"`
}

// swagger:model Experience
type Experience struct {
	// User ID
	U_ID string `json:"user_id" validate:"uuid"`
	// Experience Details
	Experience string `json:"experience"`
}

// swagger:parameters post-experience
type AddExperienceSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body Experience `json:"body"`
}

// swagger:model Skills
type Skills struct {
	// User ID
	U_ID string `json:"user_id" validate:"uuid"`
	// Experience Details
	Skill string `json:"skill"`
}

// swagger:parameters post-skills
type AddSkillSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body Skills `json:"body"`
}

type Profile struct {
	U_ID       string `json:"user_id" validate:"uuid"`
	Education  string `json:"education"`
	Phone      string `json:"phone" validate:"numeric"`
	Experience []Experience
	Skills     []Skills
}

// swagger:model Profile
type ProfileWrapper struct {
	// ID of user of which, the profile is to be added
	U_ID string `json:"user_id" validate:"uuid"`
	// user's education
	Education string `json:"education"`
	// user's Phone number
	Phone string `json:"phone" validate:"numeric"`
}

// swagger:parameters post-profile
type AddProfileSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body ProfileWrapper `json:"body"`
}

// swagger:parameters userProfile
type ProfileParam struct {
	// User ID
	// in: query
	User_ID string `json:"user_id"`
}

type Applications struct {
	Job_ID    string `json:"job_id" validate:"uuid"`
	U_ID      string `json:"user_id" validate:"uuid"`
	Status    string `json:"status"`
	CreatedAt time.Time
}

// swagger:model Applications
type ApplicationWrapper struct {
	Job_ID string `json:"job_id" validate:"uuid"`
	U_ID   string `json:"user_id" validate:"uuid"`
	Status string `json:"status"`
}

// swagger:parameters post-application
type AddApplicationSwagger struct {
	//  Add details of the user
	//  in: body
	//  required: true
	Body ApplicationWrapper `json:"body"`
}
