package integration_test

import "testing"

func TestCheckPrometheusDeveloper(t *testing.T) {
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + developer.Token). 
	Expect().
	Status(200)

}

