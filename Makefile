dev:
	npx nodemon server.js

REMOTE=ubuntu@radio.nonzerosoftware.com
deploy:
	docker build . --tag emb-radio
	docker save emb-radio:latest -o dist/emb-radio.tar
	scp -i private/deploy.rsa dist/emb-radio.tar $(REMOTE):/tmp
	-ssh -i private/deploy.rsa $(REMOTE) docker stop emb-radio
	-ssh -i private/deploy.rsa $(REMOTE) docker rm emb-radio
	ssh -i private/deploy.rsa $(REMOTE) docker load -i /tmp/emb-radio.tar
	ssh -i private/deploy.rsa $(REMOTE) docker run --detach -p3010:3010 --volume=/media:/media --name emb-radio emb-radio:latest

logs:
	ssh -i private/deploy.rsa $(REMOTE) docker logs emb-radio:latest

clean:
	rm -rf dist/*
