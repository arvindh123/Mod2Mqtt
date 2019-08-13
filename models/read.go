package models

type DeviceWithRegs struct {
	Device  DeviceDetails
	Modregs []ModbusRegisters
}
type InterWithDevices struct {
	Inter   InterfaceDetails
	Devices []DeviceWithRegs
}

type AllStructParams struct {
	AllStructParams []InterWithDevices
}
