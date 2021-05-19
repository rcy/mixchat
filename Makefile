dev:
	npx nodemon index.js

dockerize:
	docker build . --tag emb-radio
