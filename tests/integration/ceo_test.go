package integration_test

import "testing"

func TestMetricsCeo(t *testing.T) {
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + ceo.Token). 
	Expect().
	Status(200)
}

func TestUpdateBioCeo(t *testing.T) {
	body := map[string]interface{}{
		"bio": "i'm a ceo from the us",
	}
	e := newExpect(t)
	e.PATCH("/users/update/bio").
	WithHeader("Authorization", "Bearer " + ceo.Token). 
	WithJSON(body).
	Expect(). 
	Status(200)
}