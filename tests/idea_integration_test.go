package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type IdeaIntegrationTestSuite struct {
	BaseTestSuite
}

func (suite *IdeaIntegrationTestSuite) SetupSuite() {
	suite.BaseTestSuite.SetupSuite()
}

func (suite *IdeaIntegrationTestSuite) TearDownTest() {
	suite.BaseTestSuite.TearDownTest()
}

func TestIdeaIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IdeaIntegrationTestSuite))
}

// createIdeaPrerequisites is a helper to set up a user, coffee shop, and category for tests.
func (suite *IdeaIntegrationTestSuite) createIdeaPrerequisites() (string, *models.User, *models.CoffeeShop, *models.Category) {
	// Create user and get token
	user := suite.CreateUser("idea-creator", "111111111")
	token := suite.RegisterUserAndGetToken(user)

	// Create coffee shop
	coffeeShop := &models.CoffeeShop{
		Name:      "Test Coffee Shop for Ideas",
		Address:   "123 Idea St",
		CreatorID: user.ID,
	}
	err := suite.DB.Create(coffeeShop).Error
	suite.Require().NoError(err)

	// Create category
	category := &models.Category{
		Title:        "Idea Category",
		CoffeeShopID: &coffeeShop.ID,
	}
	err = suite.DB.Create(category).Error
	suite.Require().NoError(err)

	return token, user, coffeeShop, category
}

// createTestIdea is a helper to create a single idea for use in tests.
func (suite *IdeaIntegrationTestSuite) createTestIdea(author *models.User, cs *models.CoffeeShop, cat *models.Category) *models.Idea {
	var ideaStatus models.IdeaStatus
	suite.DB.FirstOrCreate(&ideaStatus, "title = ?", "new")

	idea := &models.Idea{
		Title:        "Test Idea",
		Description:  "A brilliant idea.",
		CreatorID:    &author.ID,
		CoffeeShopID: &cs.ID,
		CategoryID:   &cat.ID,
		StatusID:     &ideaStatus.ID,
	}
	err := suite.DB.Create(idea).Error
	suite.Require().NoError(err)
	return idea
}

func (suite *IdeaIntegrationTestSuite) TestCreateIdea() {
	token, _, coffeeShop, category := suite.createIdeaPrerequisites()

	// Dummy image content (minimal valid JPEG)
	dummyImageContent := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0xFF, 0xC0, 0x00, 0x11, 0x08, 0x00,
		0x01, 0x00, 0x01, 0x03, 0x01, 0x22, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11,
		0x01, 0xFF, 0xC4, 0x00, 0x1F, 0x00, 0x00, 0x01, 0x05, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0xFF, 0xC4,
		0x00, 0xB5, 0x10, 0x00, 0x02, 0x01, 0x03, 0x03, 0x02, 0x04, 0x03, 0x05,
		0x05, 0x04, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x01, 0x00, 0x03, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
		0xFF, 0xDA, 0x00, 0x0C, 0x03, 0x01, 0x00, 0x02, 0x11, 0x03, 0x11, 0x00,
		0x3F, 0x00, 0xBF, 0x40, 0xFF, 0xD9,
	}

	tests := []struct {
		name           string
		token          string
		formData       map[string]string
		fileField      string
		fileContent    []byte
		fileName       string
		contentType    string
		expectedStatus int
		checkResponse  bool
		expectedTitle  string
		expectedImageURL string
	}{
		{
			name:  "Success - Create Idea without Image",
			token: token,
			formData: map[string]string{
				"title":          "My New Awesome Idea",
				"description":    "This is the description of my awesome idea.",
				"category_id":    category.ID.String(),
				"coffee_shop_id": coffeeShop.ID.String(),
			},
			contentType:    "multipart/form-data",
			expectedStatus: http.StatusCreated,
			checkResponse:  true,
			expectedTitle:  "My New Awesome Idea",
			expectedImageURL: "",
		},
		{
			name:  "Success - Create Idea with Image",
			token: token,
			formData: map[string]string{
				"title":          "Idea with Picture",
				"description":    "This idea has a picture.",
				"category_id":    category.ID.String(),
				"coffee_shop_id": coffeeShop.ID.String(),
			},
			fileField:      "image",
			fileContent:    dummyImageContent,
			fileName:       "test_image.jpg",
			contentType:    "multipart/form-data",
			expectedStatus: http.StatusCreated,
			checkResponse:  true,
			expectedTitle:  "Idea with Picture",
			expectedImageURL: "test_image.jpg",
		},
		{
			name:           "Fail - Unauthorized",
			token:          "",
			formData:       map[string]string{},
			contentType:    "multipart/form-data",
			expectedStatus: http.StatusUnauthorized,
			expectedTitle: "",
			expectedImageURL: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := TestRequest{
				method:      http.MethodPost,
				path:        "/api/v1/ideas",
				token:       tt.token,
				formData:    tt.formData,
				fileField:   tt.fileField,
				fileContent: tt.fileContent,
				fileName:    tt.fileName,
				contentType: tt.contentType,
			}
			w := suite.MakeRequest(req)
			suite.Equal(tt.expectedStatus, w.Code)

			if tt.checkResponse {
				var resp dto.IdeaResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				suite.NoError(err)
				suite.Equal(tt.expectedTitle, resp.Title)
				suite.Equal(coffeeShop.ID, *resp.CoffeeShopID)
				suite.NotEqual(uuid.Nil, resp.ID)
				suite.False(resp.CreatedAt.IsZero(), "CreatedAt should not be zero")
				if tt.expectedImageURL != "" {
					suite.NotNil(resp.ImageURL)
					suite.Contains(*resp.ImageURL, tt.expectedImageURL)
				} else {
					suite.Nil(resp.ImageURL)
				}
			}
		})
	}
}

func (suite *IdeaIntegrationTestSuite) TestGetAllIdeas() {
	_, user, coffeeShop, category := suite.createIdeaPrerequisites()
	idea := suite.createTestIdea(user, coffeeShop, category)

	req := TestRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/api/v1/coffee-shops/%s/ideas", coffeeShop.ID),
	}
	w := suite.MakeRequest(req)

	suite.Equal(http.StatusOK, w.Code)

	var resp []dto.IdeaResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	suite.NoError(err)

	found := false
	for _, i := range resp {
		if i.ID == idea.ID {
			found = true
			suite.False(i.CreatedAt.IsZero(), "CreatedAt should not be zero")
			break
		}
	}
	suite.True(found, "created idea not found in list")
}

func (suite *IdeaIntegrationTestSuite) TestGetIdea() {
	_, user, coffeeShop, category := suite.createIdeaPrerequisites()
	idea := suite.createTestIdea(user, coffeeShop, category)

	req := TestRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
	}
	w := suite.MakeRequest(req)

	suite.Equal(http.StatusOK, w.Code)

	var resp dto.IdeaResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(idea.ID, resp.ID)
	suite.Equal(idea.Title, resp.Title)
	suite.False(resp.CreatedAt.IsZero(), "CreatedAt should not be zero")
}

func (suite *IdeaIntegrationTestSuite) TestUpdateIdea() {
	authorToken, author, coffeeShop, category := suite.createIdeaPrerequisites()
	otherToken := suite.GetRandomAuthToken()

	admin := suite.CreateUser("admin-idea", "987654321")
	adminToken := suite.RegisterUserAndGetToken(admin)
	updTitle := "Updated Title"
	updateReq := dto.UpdateIdeaRequest{Title: &updTitle}

	suite.Run("Author can update their idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method:      http.MethodPut,
			path:        fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			body:        updateReq,
			token:       authorToken,
			contentType: "application/json",
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusNoContent, w.Code)
		var updatedIdea models.Idea
		suite.DB.First(&updatedIdea, "id = ?", idea.ID)
		suite.Equal("Updated Title", updatedIdea.Title)
	})

	suite.Run("Other user cannot update idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method:      http.MethodPut,
			path:        fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			body:        updateReq,
			token:       otherToken,
			contentType: "application/json",
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusUnauthorized, w.Code)
		var notUpdatedIdea models.Idea
		suite.DB.First(&notUpdatedIdea, "id = ?", idea.ID)
		suite.Equal("Test Idea", notUpdatedIdea.Title)
	})

	suite.Run("Admin cannot update idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method:      http.MethodPut,
			path:        fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			body:        updateReq,
			token:       adminToken,
			contentType: "application/json",
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusUnauthorized, w.Code)
		var notUpdatedIdea models.Idea
		suite.DB.First(&notUpdatedIdea, "id = ?", idea.ID)
		suite.Equal("Test Idea", notUpdatedIdea.Title)
	})

	suite.Run("Unauthorized cannot update", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method:      http.MethodPut,
			path:        fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			body:        updateReq,
			token:       "",
			contentType: "application/json",
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusUnauthorized, w.Code)
	})
}

func (suite *IdeaIntegrationTestSuite) TestDeleteIdea() {
	authorToken, author, coffeeShop, category := suite.createIdeaPrerequisites()
	otherToken := suite.GetRandomAuthToken()

	admin := suite.CreateUser("admin-idea-del", "123123123")
	adminToken := suite.RegisterUserAndGetToken(admin)

	// Make the admin a worker in the shop to test admin deletion privileges
	err := suite.DB.Create(&models.WorkerCoffeeShop{
		WorkerID:     &admin.ID,
		CoffeeShopID: &coffeeShop.ID,
		RoleID:       &suite.AdminRoleID,
	}).Error
	suite.Require().NoError(err)

	suite.Run("Author can delete their idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			token:  authorToken,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusNoContent, w.Code)
		var count int64
		suite.DB.Model(&models.Idea{}).Where("id = ?", idea.ID).Count(&count)
		suite.Equal(int64(0), count)
	})

	suite.Run("Other user cannot delete idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			token:  otherToken,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusForbidden, w.Code)
		var count int64
		suite.DB.Model(&models.Idea{}).Where("id = ?", idea.ID).Count(&count)
		suite.Equal(int64(1), count)
	})

	suite.Run("Admin can delete idea", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			token:  adminToken,
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusNoContent, w.Code)
		var count int64
		suite.DB.Model(&models.Idea{}).Where("id = ?", idea.ID).Count(&count)
		suite.Equal(int64(0), count)
	})

	suite.Run("Unauthorized cannot delete", func() {
		idea := suite.createTestIdea(author, coffeeShop, category)
		req := TestRequest{
			method: http.MethodDelete,
			path:   fmt.Sprintf("/api/v1/ideas/%s", idea.ID),
			token:  "",
		}
		w := suite.MakeRequest(req)
		suite.Equal(http.StatusUnauthorized, w.Code)
	})
}

