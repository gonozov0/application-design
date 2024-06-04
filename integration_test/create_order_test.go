package integration_test

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/levigross/grequests"
)

func (suite *APITestSuite) TestCreateOrder() {
	// TODO: create user, hotels and rooms when api will be ready
	suite.T().Skip("skipping test as api is not ready")
	resp, err := grequests.Post(suite.serverURL+"orders", &grequests.RequestOptions{
		JSON: map[string]interface{}{
			"user_email": "test@test.com",
			"hotel_id":   "reddison",
			"room_id":    "lux",
			"from":       "4021-01-01T00:00:00Z",
			"to":         "4021-01-02T00:00:00Z",
		},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]uuid.UUID
	err = resp.JSON(&result)
	suite.Require().NoError(err)
	suite.Require().NotEqual(uuid.Nil, result["id"])

	// TODO: get the order and check booking when api will be ready
}
