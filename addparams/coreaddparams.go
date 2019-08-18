package addparams

import (
	"strconv"

	"github.com/arvindh123/Mod2Mqtt/models"
)

var AddFeatures map[string]interface{}
var FirstInit bool = false

func InitParams() {
	var params []models.AddonFeatures
	db := models.GetDB()
	db.Find(&params)
	AddFeatures = make(map[string]interface{}, len(params))
	for _, param := range params {
		switch param.ParamType {
		case 1: //Bool
			value, err := strconv.ParseBool(param.Value)
			if err == nil {
				AddFeatures[param.Param] = value
			}
		case 2: //UInt64
			value, err := strconv.ParseUint(param.Value, 10, 64)
			if err == nil {
				AddFeatures[param.Param] = value
			}
		case 3: // Int64
			value, err := strconv.ParseInt(param.Value, 10, 64)
			if err == nil {
				AddFeatures[param.Param] = value
			}
		case 4: // Float64
			value, err := strconv.ParseFloat(param.Value, 64)
			if err == nil {
				AddFeatures[param.Param] = value
			}
		default:
			AddFeatures[param.Param] = param.Value
		}
	}
}

func GetParams() map[string]interface{} {
	// if FirstInit == false {
	// 	InitParams()
	// 	FirstInit = true
	// }
	InitParams()
	return AddFeatures
}
