build:
	docker -H docker-engine-1:2376 build \
    	--tag gcr.io/silicon-airlock-153323/github.com/ruang-guru/grader/assignment-runner-base@sha256:8d4ac15977cf36f59a993e24f0dafbbebc3a9d808338caf53b6cac6a736f6136 "asia.gcr.io/$PROJECT_ID/$IMAGE:$GITHUB_SHA" \
    	-f ./apps/rea/grader/Dockerfile \
    	.