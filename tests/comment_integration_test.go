package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/stretchr/testify/suite"
)

type CommentIntegrationTestSuite struct {
	BaseTestSuite
}

func (suite *CommentIntegrationTestSuite) SetupSuite() {
	suite.BaseTestSuite.SetupSuite()
}

func (suite *CommentIntegrationTestSuite) TearDownTest() {
	suite.BaseTestSuite.TearDownTest()
}

func TestCommentIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CommentIntegrationTestSuite))
}

// createCommentPrerequisites creates a coffee shop, a worker, an idea, and returns necessary data.
func (suite *CommentIntegrationTestSuite) createCommentPrerequisites() (string, *models.User, *models.CoffeeShop, *models.Idea) {
	// Create user (worker) and get token
	worker := suite.CreateUser("worker", "222222222")
	token := suite.RegisterUserAndGetToken(worker)

	// Create coffee shop
	coffeeShop := &models.CoffeeShop{
		Name:      "Comment Test Shop",
		Address:   "456 Comment St",
		CreatorID: worker.ID,
	}
	err := suite.DB.Create(coffeeShop).Error
	suite.Require().NoError(err)

	// Assign user as worker (admin role by default for creator usually, but explicit here)
	workerRole := models.Role{Name: "worker"}
	suite.DB.FirstOrCreate(&workerRole, "name = ?", "worker")
	
	err = suite.DB.Create(&models.WorkerCoffeeShop{
		WorkerID:     &worker.ID,
		CoffeeShopID: &coffeeShop.ID,
		RoleID:       &workerRole.ID,
	}).Error
	suite.Require().NoError(err)

	// Create category
	category := &models.Category{
		Title:        "General",
		CoffeeShopID: &coffeeShop.ID,
	}
	err = suite.DB.Create(category).Error
	suite.Require().NoError(err)

	// Create idea status
	var ideaStatus models.IdeaStatus
	suite.DB.FirstOrCreate(&ideaStatus, "title = ?", "new")

	// Create idea
	idea := &models.Idea{
		Title:        "Idea for Comment",
		Description:  "Discuss this.",
		CreatorID:    &worker.ID,
		CoffeeShopID: &coffeeShop.ID,
		CategoryID:   &category.ID,
		StatusID:     &ideaStatus.ID,
	}
	err = suite.DB.Create(idea).Error
	suite.Require().NoError(err)

	return token, worker, coffeeShop, idea
}

func (suite *CommentIntegrationTestSuite) TestCreateComment() {
	token, _, _, idea := suite.createCommentPrerequisites()
	
	// Create another user who is NOT a worker
	outsider := suite.CreateUser("outsider", "333333333")
	outsiderToken := suite.RegisterUserAndGetToken(outsider)

	tests := []struct {
		name           string
		token          string
		body           dto.CreateCommentRequest
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:  "Success - Worker creates comment",
			token: token,
			body: dto.CreateCommentRequest{
				Text:       "Great idea!",
				AuthorName: "Barista Bob",
			},
			expectedStatus: http.StatusCreated,
			checkResponse:  true,
		},
		{
			name:  "Fail - Outsider cannot create comment",
			token: outsiderToken,
			body: dto.CreateCommentRequest{
				Text:       "I shouldn't be here",
				AuthorName: "Hacker",
			},
			expectedStatus: http.StatusForbidden,
			checkResponse:  false,
		},
		{
			name:  "Fail - Unauthorized",
			token: "",
			body: dto.CreateCommentRequest{
				Text:       "No token",
				AuthorName: "Anon",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := TestRequest{
				method:      http.MethodPost,
				path:        fmt.Sprintf("/api/v1/ideas/%s/comments", idea.ID),
				token:       tt.token,
				body:        tt.body,
				contentType: "application/json",
			}
			w := suite.MakeRequest(req)
			suite.Equal(tt.expectedStatus, w.Code)

			if tt.checkResponse {
				var resp dto.CommentResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				suite.NoError(err)
				suite.Equal(tt.body.Text, resp.Text)
				suite.Equal(tt.body.AuthorName, resp.AuthorName)
				suite.NotEmpty(resp.ID)
				suite.NotEmpty(resp.CreatedAt)
			}
		})
	}
}

func (suite *CommentIntegrationTestSuite) TestGetComments() {
	token, _, _, idea := suite.createCommentPrerequisites()

	// Seed some comments
	for i := 1; i <= 15; i++ {
		req := TestRequest{
			method:      http.MethodPost,
			path:        fmt.Sprintf("/api/v1/ideas/%s/comments", idea.ID),
			token:       token,
			body:        dto.CreateCommentRequest{Text: fmt.Sprintf("Comment %d", i), AuthorName: "Tester"},
			contentType: "application/json",
		}
		suite.MakeRequest(req)
	}

	// Create outsider
	outsider := suite.CreateUser("outsider2", "444444444")
	outsiderToken := suite.RegisterUserAndGetToken(outsider)

	suite.Run("Worker can get comments with pagination", func() {
		// Page 1, Limit 10 (should get newest 10: Comment 15 to Comment 6)
		req := TestRequest{
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments?page=1&limit=10", idea.ID),
			token:  token,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusOK, w.Code)

		var resp []dto.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		suite.NoError(err)
		suite.Len(resp, 10)
		suite.Equal("Comment 15", resp[0].Text) // Check ordering (newest first)
	})

	suite.Run("Worker can get second page", func() {
		// Page 2, Limit 10 (should get remaining 5: Comment 5 to Comment 1)
		req := TestRequest{
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments?page=2&limit=10", idea.ID),
			token:  token,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusOK, w.Code)

		var resp []dto.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		suite.NoError(err)
		suite.Len(resp, 5)
		suite.Equal("Comment 5", resp[0].Text)
	})

	suite.Run("Outsider cannot get comments", func() {
		req := TestRequest{
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments", idea.ID),
			token:  outsiderToken,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusForbidden, w.Code)
	})
}

func (suite *CommentIntegrationTestSuite) TestDeleteComment() {
	token, _, _, idea := suite.createCommentPrerequisites()

	// Create a comment
	createReq := TestRequest{
		method:      http.MethodPost,
		path:        fmt.Sprintf("/api/v1/ideas/%s/comments", idea.ID),
		token:       token,
		body:        dto.CreateCommentRequest{Text: "To be deleted", AuthorName: "Deletor"},
		contentType: "application/json",
	}
	w := suite.MakeRequest(createReq)
	var commentResp dto.CommentResponse
	json.Unmarshal(w.Body.Bytes(), &commentResp)

	// Outsider
	outsider := suite.CreateUser("outsider3", "555555555")
	outsiderToken := suite.RegisterUserAndGetToken(outsider)

	suite.Run("Outsider cannot delete comment", func() {
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments/%s", idea.ID, commentResp.ID),
			token:  outsiderToken,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusForbidden, w.Code)
	})

	suite.Run("Worker can delete comment", func() {
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments/%s", idea.ID, commentResp.ID),
			token:  token,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusNoContent, w.Code)

		// Verify deletion
		checkReq := TestRequest{
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/v1/ideas/%s/comments", idea.ID),
			token:  token,
		}
		checkW := suite.MakeRequest(checkReq)
		var listResp []dto.CommentResponse
		json.Unmarshal(checkW.Body.Bytes(), &listResp)
		
		found := false
		for _, c := range listResp {
			if c.ID == commentResp.ID {
				found = true
				break
			}
		}
		suite.False(found, "Deleted comment should not be returned")
	})
}
