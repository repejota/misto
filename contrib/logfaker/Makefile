.PHONY: docker-build
docker-build:	## Build logfaker container
	docker build --rm -t repejota/logfaker .

.PHONY: docker-run
docker-run:	## Run logfaker container
	docker run -t --name logfaker1 --rm repejota/logfaker python3 /opt/app/main.py

.PHONY: help
help:	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'