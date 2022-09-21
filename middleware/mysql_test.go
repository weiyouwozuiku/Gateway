package middleware_test

import "testing"

type demo struct {
	name     string `json:"id" `
	id       int    `json:"id" gorm:"primary_key"`
	password string `json:"password"`
}

func Test_DBPool(t *testing.T) {

}
func Test_GormPool(t *testing.T) {

}
