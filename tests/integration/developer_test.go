package integration_test

import "testing"

func TestCheckPrometheusDeveloper(t *testing.T) {
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + developer.Token). 
	Expect().
	Status(200)

}

func TestUpdateBioDev(t *testing.T) {
	body := map[string]interface{}{
		"bio": "i'm a dev from the us",
	}
	e := newExpect(t)
	e.PATCH("/users/update/bio").
	WithHeader("Authorization", "Bearer " + developer.Token). 
	WithJSON(body).
	Expect(). 
	Status(200)
}

