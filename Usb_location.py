 from serial.tools import list_ports
 k=list(list_ports.comports())
 for i in k:
	print(i.location,i.device)

