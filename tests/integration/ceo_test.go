package integration_test

import "testing"

func TestMetricsCeo(t *testing.T) {
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + ceo.Token). 
	Expect().
	Status(200)
}