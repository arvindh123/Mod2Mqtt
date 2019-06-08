from pymodbus.client.sync import ModbusTcpClient as ModbusClient

import logging
FORMAT = ('%(asctime)-15s %(threadName)-15s '
          '%(levelname)-8s %(module)-15s:%(lineno)-8s %(message)s')
logging.basicConfig(format=FORMAT)
log = logging.getLogger()
log.setLevel(logging.DEBUG)

UNIT = 0x1


def run_sync_client():
    
    client = ModbusClient('10.87.84.102', port=502)
    
    client.connect()

    log.debug("Write to a holding register and read back")

    rr = client.read_holding_registers(3019, 14, unit=UNIT)

    print(rr.registers)

    client.close()


if __name__ == "__main__":
    run_sync_client()
