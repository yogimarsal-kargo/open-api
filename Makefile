GOPRIVATE=github.com/kargotech

# Makefile guard function to ensure existence of variables
guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

ci-test: 
	go mod verify\
	&& go mod download\
	&& go test -v --race ./... -coverprofile cover.out\
	&& go tool cover -func cover.out

# Make task for CI test coverage fo CI sonarscanner
ci-test-coverage: 
	go mod verify\
	&& go mod download\
	&& go test -v --race ./... -coverprofile="coverage.out" -json > test-report.json

# Make task for CI golangci-lint report fo CI sonarscanner
ci-golangci-lint-report: 
	go mod verify\
	&& go mod download\
	&& wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.47.1\
	&& ./bin/golangci-lint run --out-format checkstyle --issues-exit-code 0 > golangci-lint.out

# Make task for CI horusec report fo CI sonarscanner
ci-horusec-report: 
	curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s v2.8.0\
	&& horusec start -p . -o="sonarqube" -O="./horusec-sonarqube.json"

# Make whole sonarqube report data from test coverage, golangci-lint, and horusec
ci-sonarqube-report: ci-test-coverage ci-golangci-lint-report ci-horusec-report

ci-package: guard-DOCKER_REGISTRY
	docker build -t $(DOCKER_REGISTRY):$(TAG) --build-arg COMMIT_HASH=$(TAG) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) . \
	&& docker push $(DOCKER_REGISTRY):$(TAG)
