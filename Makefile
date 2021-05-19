dev:
	npx nodemon server.js

dockerize:
	docker build . --tag emb-radio

REMOTE=ubuntu@radio.nonzerosoftware.com
deploy: dist/emb-radio.tar
	scp -i private/deploy.rsa $< $(REMOTE):/tmp
	ssh -i private/deploy.rsa $(REMOTE) docker run --detach emb-radio:latest

logs:
	ssh -i private/deploy.rsa $(REMOTE) docker logs emb-radio:latest

dist/emb-radio.tar: dockerize
	docker save emb-radio:latest -o $@

clean:
	rm -rf dist/*
