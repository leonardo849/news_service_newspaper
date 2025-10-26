package integration_test

import (
	"news_service/internal/dto"
	"testing"
)

var id  string

func TestCreateNews(t *testing.T) {
	news := map[string]interface{}{
		"title":    "the amazing world of gumball",
		"subtitle": "a cartoon about a crazy blue cat",
		"topic":    "entertainment",
		"blocks": []map[string]interface{}{
			{
				"content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
				"position": 1,
				"images": []map[string]string{
					{
						"url": "https://static.wikia.nocookie.net/c29c5a0d-cbca-4ac3-ab27-63be0b47bae0/scale-to-width/755",
					},
				},
			},
		},
	}
	e := newExpect(t)
	response := e.POST("/news/create").
	WithHeader("Authorization", "Bearer " + journalist.Token).
	WithJSON(news). 
	Expect(). 
	Status(200).
	JSON().Object()

	id = response.Value("id").String().Raw()
}

func TestCreateBlocks(t *testing.T) {
	blocks := []dto.CreateBlockDTO{
		{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
			Position: 2,
			Images: nil,
		},
		{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
			Position: 3,
			Images: nil,
		},
	}
	e := newExpect(t)
	e.POST("/blocks/create/" + id). 
	WithHeader("Authorization", "Bearer " + journalist.Token).
	WithJSON(blocks). 
	Expect().Status(200)
}

func TestPublishNews(t *testing.T) {
	
	e := newExpect(t)
	e.PATCH("/news/update/publish/" + id).
	WithHeader("Authorization", "Bearer " + journalist.Token).
	Expect(). 
	Status(200)
}

func TestUpdateBio(t *testing.T) {
	body := map[string]interface{}{
		"bio": "i'm a journalist from the us",
	}
	e := newExpect(t)
	e.PATCH("/users/update/bio").
	WithHeader("Authorization", "Bearer " + journalist.Token). 
	WithJSON(body).
	Expect(). 
	Status(200)
}