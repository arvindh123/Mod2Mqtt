package models

type AddonFeatures struct {
	ID        int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Param     string `gorm:"UNIQUE;NOT NULL" json:"param"`
	Value     string `gorm:"NOT NULL" json:"value"`
	ParamType int    `gorm:"NOT NULL" json:"paramtype"`
}
