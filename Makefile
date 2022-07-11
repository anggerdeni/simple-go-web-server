GITHUB_SHA=$(shell git rev-parse --verify HEAD)
PROJECT_ID=silicon-airlock-153323
IMAGE=github.com/ruang-guru/grader/assignment-runner-base

build:
	docker build --tag gcr.io/${PROJECT_ID}/${IMAGE}:${GITHUB_SHA} .

publish:
	docker push gcr.io/${PROJECT_ID}/${IMAGE}:${GITHUB_SHA}