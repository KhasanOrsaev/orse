MIGRATION_DIR = "var/migration/"

.PHONY: migrate
migrate: deps.migrate
	MIGRATION_VERSION=0
#ifneq ($(wildcard ${MIGRATION_DIR}version.txt),) 
#	MIGRATION_VERSION := $(shell cat ${MIGRATION_DIR}version.txt)
#endif
	@m=333
	@for file in $(MIGRATION_DIR)/*.sql; do \
		echo $$m; \
		VERSION=$$(echo $$file | sed -E 's/.*\/([0-9]+)_.*/\1/g'); \
		if [ ${MIGRATION_VERSION} -lt $$VERSION ]; then \
		sqlite3 home.db < var/migration/20221230_init.sql; \
		MIGRATION_VERSION = $$VERSION; \
		echo ${MIGRATION_VERSION} > ${MIGRATION_DIR}version.txt; \
		fi \
	done



deps.migrate:
ifeq (, $(shell which sqlite3)) 
	echo "sqlite3 isn't installed. please install sqlite3 client"
endif
