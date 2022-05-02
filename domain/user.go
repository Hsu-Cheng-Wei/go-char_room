package domain

type User struct {
	ID       string `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null;" json:"name"`
	Email    string `gorm:"type:varchar(20);not null;" json:"email"`
	Password string `gorm:"type:varchar(46);not null;" json:"password"`
}
