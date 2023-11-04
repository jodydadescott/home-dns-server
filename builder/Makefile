default:
	@echo "You really want to build all or maybe chdir into what you want? Or use make all"

all:
	cd darwin && $(MAKE)
	cd linux && $(MAKE)

clean:
	cd darwin && $(MAKE) clean
	cd linux && $(MAKE) clean

push:
	cd linux && $(MAKE) push