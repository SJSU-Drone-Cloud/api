package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tracking struct {
	ID           primitive.ObjectID  `json:"_id" bson:"_id"`
	DroneID      string              `json:"droneID" bson:"droneID"`
	TimeLocation []TimeLocationStamp `json:"timeLocation" bson:"timeLocation"`
	LastUpdated  time.Time           `json:"lastUpdated" bson:"lastUpdated"`
}

type TrackingDevice struct {
	DroneID string  `json:"droneID" bson:"droneID"`
	Lat     float64 `json:"lat,string" bson:"lat"`
	Lng     float64 `json:"lng,string" bson:"lng"`
}

type TimeLocationStamp struct {
	TimeStamp time.Time `json:"timestamp" bson:"timestamp"`
	Lat       float64   `json:"lat" bson:"lat"`
	Lng       float64   `json:"lng" bson:"lng"`
}

type Drone struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	DroneID     string             `json:"droneID" bson:"droneID"`
	Coordinates Coordinates        `json:"coordinates" bson:"coordinates"`
	Address     string             `json:"address" bson:"address"`
	LastUpdated time.Time          `json:"lastUpdated" bson:"lastUpdated"`
}

type RegisterDrone struct {
	Address string  `json:"address" bson:"address"`
	Lat     float64 `json:"lat" bson:"lat"`
	Lng     float64 `json:"lng" bson:"lng"`
}

type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

func (rd *RegisterDrone) UnmarshalJSON(data []byte) error {
	var a map[string]string
	fmt.Println("registerDrone")
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom regdrone unmarshal error")
		return err
	}
	rd.Address = a["address"]
	lat, err := strconv.ParseFloat(a["lat"], 64)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(a["lng"], 64)
	if err != nil {
		return err
	}
	rd.Lat = lat
	rd.Lng = lng
	return nil
}
