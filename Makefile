

APP = app
GOFILES = $(wildcard ./src/*.go)
TMP_FILES = $(wildcard ./*~) $(wildcard ./src/*~)

default: $(APP)

$(APP): $(GOFILES)
	go build -o $@ $^ && strip $(APP)

clean:
	-rm -rf $(APP) $(TMP_FILES)

