package main

import "errors"

type Receiver struct {
	Name     string `json:"name" bson:"name"`
	Type     uint8  `json:"type" bson:"type"`
	Endpoint string `json:"endpoint" bson:"endpoint"`
}

// - Repository

func getReceivers(userID string) (*[]Receiver, error) {
	u, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &u.Receivers, nil
}

func (r Receiver) createReceiver(userID string) error {
	u, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	u.Receivers = append(u.Receivers, r)

	return users.UpdateId(userID, &u)
}

func removeReceiver(userID string, name string) error {
	u, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	for i, rec := range u.Receivers {
		if rec.Name == name {
			u.Receivers = append(u.Receivers[:i], u.Receivers[i+1:]...)
			return users.UpdateId(userID, &u)
		}
	}

	return errors.New("Receiver " + name + " for user " + u.Fullname + " not found")
}
