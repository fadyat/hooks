PORT=80

ifneq (,$(wildcard ./.env))
	include .env
	export
endif