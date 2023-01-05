//go:build all
// +build all

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Massad/gin-boilerplate/controllers"
	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

var auth = new(controllers.AuthController)

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		v1.POST("/token/refresh", auth.Refresh)

		/*** START Product ***/
		product := new(controllers.ProductController)

		v1.POST("/product", TokenAuthMiddleware(), product.Create)
		v1.GET("/products", TokenAuthMiddleware(), product.All)
		v1.GET("/product/:id", TokenAuthMiddleware(), product.One)
		v1.PUT("/product/:id", TokenAuthMiddleware(), product.Update)
		v1.DELETE("/product/:id", TokenAuthMiddleware(), product.Delete)
	}

	return r
}

func main() {
	r := SetupRouter()
	r.Run()
}

var loginCookie string

var testEmail = "test-gin-boilerplate@test.com"
var testPassword = "123456"

var accessToken string
var refreshToken string

var productID int

/**
* TestIntDB
* It tests the connection to the database and init the db for this test
*
* Must pass
 */
func TestIntDB(t *testing.T) {

	//Load the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}

	fmt.Println("DB_PASS", os.Getenv("DB_PASS"))

	db.Init()
	db.InitRedis(1)
}

/**
* TestRegister
* Test user registration
*
* Must return response code 200
 */
func TestRegister(t *testing.T) {
	testRouter := SetupRouter()

	var registerForm forms.RegisterForm

	registerForm.Name = "testing"
	registerForm.Email = testEmail
	registerForm.Password = testPassword

	data, _ := json.Marshal(registerForm)

	req, err := http.NewRequest("POST", "/v1/user/register", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestRegisterInvalidEmail
* Test user registration with invalid email
*
* Must return response code 406
 */
func TestRegisterInvalidEmail(t *testing.T) {
	testRouter := SetupRouter()

	var registerForm forms.RegisterForm

	registerForm.Name = "testing"
	registerForm.Email = "invalid@email"
	registerForm.Password = testPassword

	data, _ := json.Marshal(registerForm)

	req, err := http.NewRequest("POST", "/v1/user/register", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}

/**
* TestLogin
* Test user login
* and get the access_token and refresh_token stored
*
* Must return response code 200
 */
func TestLogin(t *testing.T) {
	testRouter := SetupRouter()

	var loginForm forms.LoginForm

	loginForm.Email = testEmail
	loginForm.Password = testPassword

	data, _ := json.Marshal(loginForm)

	req, err := http.NewRequest("POST", "/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		Message string `json:"message"`
		User    struct {
			CreatedAt int64  `json:"created_at"`
			Email     string `json:"email"`
			ID        int64  `json:"id"`
			Name      string `json:"name"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"user"`
		Token struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"token"`
	}
	json.Unmarshal(body, &res)

	accessToken = res.Token.AccessToken
	refreshToken = res.Token.RefreshToken

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestInvalidLogin
* Test invalid login
*
* Must return response code 406
 */
func TestInvalidLogin(t *testing.T) {
	testRouter := SetupRouter()

	var loginForm forms.LoginForm

	loginForm.Email = "wrong@email.com"
	loginForm.Password = testPassword

	data, _ := json.Marshal(loginForm)

	req, err := http.NewRequest("POST", "/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}

/**
* TestCreateProduct
* Test product creation
*
* Must return response code 200
 */
func TestCreateProduct(t *testing.T) {
	testRouter := SetupRouter()

	var form forms.CreateProductForm

	form.Title = "Testing product title"
	form.Content = "Testing product content"

	data, _ := json.Marshal(form)

	req, err := http.NewRequest("POST", "/v1/product", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		Status int
		ID     int
	}
	json.Unmarshal(body, &res)

	productID = res.ID

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestCreateInvalidProduct
* Test product invalid creation
*
* Must return response code 406
 */
func TestCreateInvalidProduct(t *testing.T) {
	testRouter := SetupRouter()

	var form forms.CreateProductForm

	form.Title = "Testing product title"

	data, _ := json.Marshal(form)

	req, err := http.NewRequest("POST", "/v1/product", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}

/**
* TestGetProduct
* Test getting one product
*
* Must return response code 200
 */
func TestGetProduct(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/product/%d", productID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestGetInvalidProduct
* Test getting invalid product
*
* Must return response code 404
 */
func TestGetInvalidProduct(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/v1/product/invalid", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

/**
* TestGetProductNotLoggedin
* Test getting the product with logged out user
*
* Must return response code 401
 */
func TestGetProductNotLoggedin(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/product/%d", productID), nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

/**
* TestGetProductUnauthorized
* Test getting the product with unauthorized user (wrong or expired access_token)
*
* Must return response code 401
 */
func TestGetProductUnauthorized(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/product/%d", productID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", "abc123"))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

/**
* TestUpdateProduct
* Test updating an product
*
* Must return response code 200
 */
func TestUpdateProduct(t *testing.T) {
	testRouter := SetupRouter()

	var form forms.CreateProductForm

	form.Title = "Testing new product title"
	form.Content = "Testing new product content"

	data, _ := json.Marshal(form)

	url := fmt.Sprintf("/v1/product/%d", product)

	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestDeleteProduct
* Test deleting an product
*
* Must return response code 200
 */
func TestDeleteProduct(t *testing.T) {
	testRouter := SetupRouter()

	url := fmt.Sprintf("/v1/product/%d", productID)

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestRefreshToken
* Test refreshing the token with valid refresh_token
*
* Must return response code 200
 */
func TestRefreshToken(t *testing.T) {
	testRouter := SetupRouter()

	var tokenForm forms.Token

	tokenForm.RefreshToken = refreshToken

	data, _ := json.Marshal(tokenForm)

	req, err := http.NewRequest("POST", "/v1/token/refresh", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestInvalidRefreshToken
* Test refreshing the token with invalid refresh_token
*
* Must return response code 401
 */
func TestInvalidRefreshToken(t *testing.T) {
	testRouter := SetupRouter()

	var tokenForm forms.Token

	//Since we didn't update it in the test before - this will not be valid anymore
	tokenForm.RefreshToken = refreshToken

	data, _ := json.Marshal(tokenForm)

	req, err := http.NewRequest("POST", "/v1/token/refresh", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

/**
* TestUserSignout
* Test logout a user
*
* Must return response code 200
 */
func TestUserLogout(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/v1/user/logout", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

/**
* TestCleanUp
* Deletes the created user with it's products
*
* Must pass
 */
func TestCleanUp(t *testing.T) {
	var err error
	_, err = db.GetDB().Exec("DELETE FROM public.user WHERE email=$1", testEmail)
	if err != nil {
		t.Error(err)
	}
}
