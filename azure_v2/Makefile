#! /usr/bin/make

NAME=azure_v2
# the default target builds a binary in the top-level dir for whatever the local OS is
default: $(NAME)
$(NAME): *.go
	go build -o $(NAME) .

# the standard build produces a "local" executable, a linux tgz, and a darwin (macos) tgz
build: $(NAME) binary/$(NAME)-linux-amd64.tgz binary/$(NAME)-darwin-amd64.tgz

# create a tgz with the binary and any artifacts that are necessary
# note the hack to allow for various GOOS & GOARCH combos, sigh
binary/$(NAME)-%.tgz: *.go
	rm -rf binary/$(NAME)
	mkdir -p binary/$(NAME)
	tgt=$*; GOOS=$${tgt%-*} GOARCH=$${tgt#*-} go build -o binary/$(NAME)/$(NAME) .
	chmod +x binary/$(NAME)/$(NAME)
	tar -zcf $@ -C binary ./$(NAME)
	rm -r binary/$(NAME)