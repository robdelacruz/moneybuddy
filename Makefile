# Usage:
# 'make dep' and 'make webtools' to install dependencies.
# 'make clean' to clear all work files
# 'make' to build css and js into static/
# 'make serve' to start dev webserver

JSFILES = index.js helpers.js data.js
JSFILES2 = Index.svelte Accounts.svelte Txns.svelte

SRCS = t.go util.go web.go
SRCS2 = db.go dbcurrency.go dbaccount.go dbtxn.go

all: t static/style.css static/bundle.js

dep:
	sudo apt update
	sudo apt install curl software-properties-common
	curl -sL https://deb.nodesource.com/setup_13.x | sudo bash -
	sudo apt install nodejs
	sudo npm --force install -g npx
	go get -u github.com/mattn/go-sqlite3
	go get -u golang.org/x/crypto/bcrypt

webtools:
	npm install --save-dev tailwindcss
	npm install --save-dev postcss-cli
	npm install --save-dev cssnano
	npm install --save-dev svelte
	npm install --save-dev rollup
	npm install --save-dev rollup-plugin-svelte
	npm install --save-dev @rollup/plugin-node-resolve

static/style.css: twsrc.css
	#npx tailwind build twsrc.css -o twsrc.o 1>/dev/null
	#npx postcss twsrc.o > static/style.css
	npx tailwind -i twsrc.css -o static/style.css 1>/dev/null

static/bundle.js: $(JSFILES) $(JSFILES2)
	npx rollup -c

t: $(SRCS) $(SRCS2)
	go build -o t $(SRCS) $(SRCS2) $(SRCS3)

clean:
	rm -rf t static/*.js static/*.css static/*.map

serve:
	python -m SimpleHTTPServer

