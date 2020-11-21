package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func TestGetAllCustomer(t *testing.T) {
	seedCustomer(t)
	var custs []Customer

	res := request(http.MethodGet, uri("customers"), nil)
	err := res.Decode(&custs)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(custs), 0)
}

func TestCreateCustomer(t *testing.T) {
	body := bytes.NewBufferString(`{
		"name": "anuchito",
		"email": "anuchito@imail.com",
		"status": "active"
	}`)
	var cust Customer

	res := request(http.MethodPost, uri("customers"), body)
	err := res.Decode(&cust)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, cust.ID)
	assert.Equal(t, "anuchito", cust.Name)
	assert.Equal(t, "anuchito@imail.com", cust.Email)
	assert.Equal(t, "active", cust.Status)
}

func TestGetCustomerByID(t *testing.T) {
	c := seedCustomer(t)

	var lastCust Customer
	res := request(http.MethodGet, uri("customers", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&lastCust)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, lastCust.ID)
	assert.NotEmpty(t, lastCust.Name)
	assert.NotEmpty(t, lastCust.Email)
	assert.NotEmpty(t, lastCust.Status)
}

func TestUpdateCustomer(t *testing.T) {
	c := Customer{
		Name:   "nong",
		Email:  "nong@imail.com",
		Status: "inactive",
	}
	payload, _ := json.Marshal(c)
	res := request(http.MethodPut, uri("customers", strconv.Itoa(seedCustomer(t).ID)), bytes.NewBuffer(payload))
	var info Customer
	err := res.Decode(&info)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	// assert.Equal(t, c.ID, info.ID)
	assert.Equal(t, c.Name, info.Name)
	assert.Equal(t, c.Email, info.Email)
	assert.Equal(t, c.Status, info.Status)
}

func seedCustomer(t *testing.T) Customer {
	var c Customer
	body := bytes.NewBufferString(`{
		"name": "anuchito",
		"email": "anuchito@imail.com",
		"status": "active"
	}`)
	err := request(http.MethodPost, uri("customers"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create customer:", err)
	}
	return c
}

func TestDeleteCustomerByID(t *testing.T) {

	var payload map[string]string
	res := request(http.MethodDelete, uri("customers", strconv.Itoa(seedCustomer(t).ID)), nil)
	err := res.Decode(&payload)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, map[string]string{"message": "customer deleted"}, payload)
}

func uri(paths ...string) string {
	host := "http://localhost:2009"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "November 10, 2009")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
