select * from topics  
 select * from modbus_registers 
 select * from serial_details 
	select * from regs_topics 
	select * from regs_ports
	
	
	select distinct SD.id,SD.Com_port,MR.id,MR.Name FROM topics TT 
			INNER join regs_topics RT ON TT.id = RT.topics_id 
			INNER join modbus_registers MR ON MR.id = RT.modbus_registers_id 
			INNER Join regs_ports  RP ON RP.modbus_registers_id = MR.id
			inner join serial_details SD ON SD.id = RP.serial_details_id
		where TT.id = 1 AND SD.id IN (10,11)
		ORDER BY SD.ID
	
	