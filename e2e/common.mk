include ../e2e.env

test:
	go test -count=1 -timeout=10m -v a_test.go

update-golden:
	E2E_UPDATE_GOLDEN=true $(MAKE) test
