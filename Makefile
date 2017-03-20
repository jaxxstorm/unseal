clean:
	@rm -rf ./dist

build: clean
	@goxc -pv=v$(version)

version:
	@echo $(RELEASE)
