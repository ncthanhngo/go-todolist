package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Image struct {
	Id        int    `json:"id" gorm:"column:id;"`
	Url       string `json:"Url" gorm:"column:url"`
	Width     int    `json:"Width" gorm:"column:width"`
	Height    int    `json:"height" gorm:"column:height"`
	CloudName string `json:"cloud_name,omitempty" gorm:"-"`
	Extension string `json:"extension,omitempty" gorm:"-"`
}

func (Image) TableName() string {
	return "images"
}
func (j *Image) Fullfill(domain string) {
	j.Url = fmt.Sprintf("%s/%s", domain, j.Url)
}

// Lay du lieu tu DB
func (j *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal data from DB: ", value))
	}
	var img Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}
	*j = img // chuyen con tro ve day khi img da co gia tri
	return nil
}

// value return json value, implement driver.value interface
// Chuyen du lieu xuong DB
func (j *Image) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
