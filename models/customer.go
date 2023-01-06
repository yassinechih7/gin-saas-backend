package models

import (
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
)

// Customer ...
type Customer struct {
	ID        int64    `db:"id, primarykey, autoincrement" json:"id"`
	UserID    int64    `db:"user_id" json:"-"`
	Title     string   `db:"title" json:"title"`
	Content   string   `db:"content" json:"content"`
	UpdatedAt int64    `db:"updated_at" json:"updated_at"`
	CreatedAt int64    `db:"created_at" json:"created_at"`
	User      *JSONRaw `db:"user" json:"user"`
}

// CustomerModel ...
type CustomerModel struct{}

// Create ...
func (m CustomerModel) Create(userID int64, form forms.CreateCustomerForm) (customerID int64, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.customer(user_id, title, content) VALUES($1, $2, $3) RETURNING id", userID, form.Title, form.Content).Scan(&customerID)
	return customerID, err
}

// One ...
func (m CustomerModel) One(userID, id int64) (customer Customer, err error) {
	err = db.GetDB().SelectOne(&customer, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.customer a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 AND a.id=$2 LIMIT 1", userID, id)
	return customer, err
}

// All ...
func (m CustomerModel) All(userID int64) (customers []DataList, err error) {
	_, err = db.GetDB().Select(&customers, "SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') AS data, (SELECT row_to_json(n) FROM ( SELECT count(a.id) AS total FROM public.customer AS a WHERE a.user_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.customer a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 ORDER BY a.id DESC) d", userID)
	return customers, err
}

// Update ...
func (m CustomerModel) Update(userID int64, id int64, form forms.CreateCustomerForm) (err error) {
	//METHOD 1
	//Check the customer by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	operation, err := db.GetDB().Exec("UPDATE public.customer SET title=$2, content=$3 WHERE id=$1", id, form.Title, form.Content)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("updated 0 records")
	}

	return err
}

// Delete ...
func (m CustomerModel) Delete(userID, id int64) (err error) {

	operation, err := db.GetDB().Exec("DELETE FROM public.customer WHERE id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
