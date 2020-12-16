.PHONY: test clean

test:
	go test -v

clean:
	rm -rf ./log
	rm tmp.log