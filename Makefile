compile:
	pio run
upload:
	pio run --target upload
listen:
	picocom -b 115200 /dev/ttyUSB0
.PHONY: compile upload listen
