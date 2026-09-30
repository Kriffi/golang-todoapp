include .env
export


export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume файла окружение? Опастность утери данных. [y/n]: " ans; \
	if [ "$$ans" = "y" ] ; then \
		docker compose down todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Файлы окружение очищены"; \
	else \
		echo "Очистка окружения отменена" ; \
	fi

