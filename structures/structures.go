package structures

type Users struct {
	ID    string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name  string `json:"name" validate:"alpha"`
	Email string `json:"email" validate:"email"`
	Pass  string `json:"pass" validate:"alphanum"`
}

type Organizations struct {
	Org_ID  string `json:"id" validate:"uuid" gorm:"primaryKey;autoIncrement:false"`
	Name    string `json:"name"`
	About   string `json:"about"`
	Website string `json:"website"`
	U_ID    string `json:"user_id" validate:"uuid"`
	Users   Users  `gorm:"foreignKey:U_ID;references:ID"`
}
