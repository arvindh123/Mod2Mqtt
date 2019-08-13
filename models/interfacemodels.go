package models

import "errors"

type InterfaceDetails struct {
	ID       int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"UNIQUE;NOT NULL" json:"name"`
	Type     int    `gorm:"default:0" json:"type"`
	IpAdd    string `gorm:"default:''" json:"ipadd"`
	Port     int    `gorm:"default:0" json:"port"`
	ComPort  string `gorm:"default:''" json:"comport"`
	BaudRate int    `gorm:"default:0" json:"baudrate"`
	DataBits int    `gorm:"default:0" json:"databits"`
	Parity   string `gorm:"default:''" gorm:"size:1" json:"parity"`
	StopBits int    `gorm:"default:0" json:"stopbits"`
	Timeout  int    `gorm:"default:0" json:"timeout"`
	DaqRate  int    `gorm:"default:0" json:"daqrate"`
}

func (inter *InterfaceDetails) BeforeDelete() (err error) {
	var devices []DeviceDetails
	if err := db.Where("interface_details_id = ?", inter.ID).Find(&devices).Error; err != nil {
		if err.Error() != "record not found" {
			return nil
		}
		return err
	}
	if len(devices) <= 0 {
		return nil
	}
	return errors.New("Model Can't delete until related Device removed")

}
