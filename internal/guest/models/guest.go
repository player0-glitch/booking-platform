package guest

import "gorm.io/gorm"

/*
Guest that only visis. Holds credential that can turn this
prospective guest to auser
*/
type Guest struct {
	gorm.Model
}
