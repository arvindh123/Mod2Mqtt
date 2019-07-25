	select distinct SD.id FROM topics TT 
			INNER join regs_topics RT ON TT.id = RT.topics_id 
			INNER join modbus_registers MR ON MR.id = RT.modbus_registers_id 
			INNER Join regs_ports  RP ON RP.modbus_registers_id = MR.id
			inner join serial_details SD ON SD.id = RP.serial_details_id
		