package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CategoryIntegrationTestSuite struct {
	BaseTestSuite
}

func TestCategoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryIntegrationTestSuite))
}

func (suite *CategoryIntegrationTestSuite) TestCategoryCRUD() {
	t := suite.T()

	// Create admin and coffee shop
	adminName := "Test Admin"
	adminPhone := "1111111111"
	admin := models.User{
		ID:    uuid.New(),
		Name:  &adminName,
		Phone: &adminPhone,
	}
	err := suite.DB.Create(&admin).Error
	suite.Require().NoError(err)


	adminToken := suite.RegisterUserAndGetToken(&admin)
	suite.Require().NotEmpty(adminToken)

	coffeeShop := models.CoffeeShop{
		ID:      uuid.New(),
		Name:    "Test Coffee Shop",
		Address: "123 Test St",
		CreatorID: admin.ID,
	}
	err = suite.DB.Create(&coffeeShop).Error
	suite.Require().NoError(err)

	worker := models.WorkerCoffeeShop{
		WorkerID:     &admin.ID,
		CoffeeShopID: &coffeeShop.ID,
		RoleID:       &suite.AdminRoleID,
	}
	err = suite.DB.Create(&worker).Error
	suite.Require().NoError(err)

	var createdCategoryID uuid.UUID

	t.Run("Create Category", func(t *testing.T) {
		createCategoryDTO := dto.CreateCategory{
			Title:       "New Category",
			Description: new(string),
		}
		*createCategoryDTO.Description = "A description for the new category"

		body, err := json.Marshal(createCategoryDTO)
		suite.Require().NoError(err)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/coffee-shops/%s/categories", coffeeShop.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		rr := httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusCreated, rr.Code)

		var response map[string]string
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(response["id"])

		createdCategoryID, err = uuid.Parse(response["id"])
		suite.Require().NoError(err)
	})

	t.Run("Get Category By ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/coffee-shops/%s/categories/%s", coffeeShop.ID, createdCategoryID), nil)
		
		rr := httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusOK, rr.Code)

		var category dto.CategoryResponse
		err := json.Unmarshal(rr.Body.Bytes(), &category)
		suite.Require().NoError(err)
		suite.Require().Equal("New Category", category.Title)
		suite.Require().Equal("A description for the new category", *category.Description)
	})

	t.Run("Get Categories By Coffee Shop", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/coffee-shops/%s/categories", coffeeShop.ID), nil)

		rr := httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusOK, rr.Code)

		var response struct {
			Data  []dto.CategoryResponse `json:"data"`
			Total int                    `json:"total"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		suite.Require().NoError(err)
		suite.Require().Equal(1, response.Total)
		suite.Require().Equal("New Category", response.Data[0].Title)
	})

	t.Run("Update Category", func(t *testing.T) {
		updateCategoryDTO := dto.UpdateCategory{
			Title:       "Updated Category",
			Description: new(string),
		}
		*updateCategoryDTO.Description = "An updated description"

		body, err := json.Marshal(updateCategoryDTO)
		suite.Require().NoError(err)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/coffee-shops/%s/categories/%s", coffeeShop.ID, createdCategoryID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		rr := httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusOK, rr.Code)

		// Verify update
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/coffee-shops/%s/categories/%s", coffeeShop.ID, createdCategoryID), nil)
		rr = httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		var category dto.CategoryResponse
		err = json.Unmarshal(rr.Body.Bytes(), &category)
		suite.Require().NoError(err)
		suite.Require().Equal("Updated Category", category.Title)
		suite.Require().Equal("An updated description", *category.Description)
	})

	t.Run("Delete Category", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/coffee-shops/%s/categories/%s", coffeeShop.ID, createdCategoryID), nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		rr := httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusOK, rr.Code)

		// Verify delete
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/coffee-shops/%s/categories/%s", coffeeShop.ID, createdCategoryID), nil)
		rr = httptest.NewRecorder()
		suite.Router.ServeHTTP(rr, req)

		suite.Require().Equal(http.StatusNotFound, rr.Code)
	})
}