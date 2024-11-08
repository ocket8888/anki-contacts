// Package utils contains most of the functional logic for anki-contacts,
// including handlers and database functions.
package utils

import (
	"fmt"
	"os"
	"time"

	"anki-contacts/databased"
)

// A Contact models a single contact.
type Contact struct {
	// Primary key.
	ID uint `gorm:"primaryKey"`
	// The contact's given name.
	FirstName string `gorm:"type:varchar(100);not null"`
	// The contact's family name.
	LastName string `gorm:"type:varchar(100);not null"`
	// The contact's email address.
	Email *string `gorm:"type:varchar(100);unique"`
	// The contact's telephone number.
	PhoneNumber *string `gorm:"type:varchar(20)"`
	// The contact's birthday (in UTC - hopefully).
	Birthday *time.Time
	// Automatically managed by GORM for creation time.
	CreatedAt time.Time
	// Automatically managed by GORM for last contact time.
	LastContacted time.Time
}

// MigrateSchema sets up the database schema necessary for the application to
// function.
func MigrateSchema(db databased.Database) {
	err := db.AutoMigrate(Contact{})
	if err != nil {
		panic(err)
	}
}

// CreateContact creates a new contact.
func CreateContact(contact *Contact) error {
	db := databased.Get()
	err := db.Create(contact).Error
	return err
}

// ReadAllContacts retrieves all contacts from the database.
func ReadAllContacts() ([]Contact, error) {
	db := databased.Get()
	var contacts []Contact
	tx := db.Begin()
	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Fprintf(os.Stderr, "rolling back transaction: %v\n", err)
		}
	}()
	tx.Find(&contacts)

	return contacts, tx.Error
}

// ReadContactByID retrieves a contact by uint ID
func ReadContactByID(ID uint) (Contact, error) {
	db := databased.Get()
	var contact Contact
	tx := db.Begin()
	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Fprintf(os.Stderr, "rolling back transaction: %v\n", err)
		}
	}()
	tx.Where("id = ?", ID).Find(&contact)
	return contact, tx.Error
}

// UpdateContactByID updates a contact definition.
func (contact *Contact) UpdateContactByID(contactUpdate Contact) error {
	db := databased.Get()
	err := db.Model(contact).Updates(contactUpdate).Error
	return err
}

// DeleteContact deletes a contact by uint ID.
func DeleteContactByID(ID uint) error {
	db := databased.Get()
	err := db.Where("id = ?", ID).Delete(Contact{}).Error
	return err
}
