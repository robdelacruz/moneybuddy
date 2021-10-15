## Money Buddy

Web-based app for keeping track of your money. Manage bank accounts, stocks.

Uses the following web tools:
- [Tailwind CSS](https://tailwindcss.com)
- [Svelte](https://svelte.dev)
- [Rollup](https://rollupjs.org)

`postcss` and `cssnano` are used to optionally compress css if needed.

`Make` is used to build rather than *npm*. No *package.json* file is used.

## Usage

Run once:

    $ make dep
    $ make depgo
    $ make webtools

Build and test:

    $ make clean
    $ make

See the comments in *Makefile*.

Create a new database file:

    $ mb -i myaccounts.db

Start webserver:

    $ mb myaccounts.db


## Contact
    Email: robdelacruz@protonmail.com
    Source: http://github.com/robdelacruz/moneybuddy

