package entity

import "time"

type Hospitals []*Hospital

type Hospital struct {
	Name         string    `csv:"医療機関名"`
	PhoneNumber  string    `csv:"電話番号"`
	Address      string    `csv:"正規化住所"`
	PostalCode   string    `csv:"郵便番号"`
	Municipality string    `csv:"区市町村"`
	CreatedAt    time.Time `gorm:"autoCreateTime:nano"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:nano"`
}
