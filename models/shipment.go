package models

import (
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
)

// Shipment ...
type Shipment struct {
	ID        int64    `db:"id, primarykey, autoincrement" json:"id"`
	UserID    int64    `db:"user_id" json:"-"`
	Title     string   `db:"title" json:"title"`
	Content   string   `db:"content" json:"content"`
	UpdatedAt int64    `db:"updated_at" json:"updated_at"`
	CreatedAt int64    `db:"created_at" json:"created_at"`
	User      *JSONRaw `db:"user" json:"user"`
}

// ShipmentModel ...
type ShipmentModel struct{}

// Create ...
func (m ShipmentModel) Create(userID int64, form forms.CreateShipmentForm) (shipmentID int64, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.shipment(user_id, title, content) VALUES($1, $2, $3) RETURNING id", userID, form.Title, form.Content).Scan(&shipmentID)
	return shipmentID, err
}

// One ...
func (m ShipmentModel) One(userID, id int64) (shipment Shipment, err error) {
	err = db.GetDB().SelectOne(&shipment, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.shipment a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 AND a.id=$2 LIMIT 1", userID, id)
	return shipment, err
}

// All ...
func (m ShipmentModel) All(userID int64) (shipments []DataList, err error) {
	_, err = db.GetDB().Select(&shipments, "SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') AS data, (SELECT row_to_json(n) FROM ( SELECT count(a.id) AS total FROM public.shipment AS a WHERE a.user_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.shipment a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 ORDER BY a.id DESC) d", userID)
	return shipments, err
}

// Update ...
func (m ShipmentModel) Update(userID int64, id int64, form forms.CreateShipmentForm) (err error) {
	//METHOD 1
	//Check the shipment by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	operation, err := db.GetDB().Exec("UPDATE public.shipment SET title=$2, content=$3 WHERE id=$1", id, form.Title, form.Content)
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
func (m ShipmentModel) Delete(userID, id int64) (err error) {

	operation, err := db.GetDB().Exec("DELETE FROM public.shipment WHERE id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
