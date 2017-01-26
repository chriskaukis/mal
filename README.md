# MAL
This is a Go interface to the [My Anime List API](https://myanimelist.net/modules.php?go=api)

## Getting Started
The best way to get started would be to look at the `mal_test.go` file to understand the different endpoints and what they return until I take a moment and write more documentation. The My Anime List API doesn't have many endpoints, so it's not overly complicated and easy to understand quickly. There are no 3rd party libraries being used, so it will run out of the box as they say.

## Running the tests
You will need to provide a username and password. Right now, by default it looks for the following environment variables:
`MAL_USERNAME` and `MAL_PASSWORD`

which is your [myanimelist.net username](https://myanimelist.net) and
[myanimelist.net password](https://myanimelist.net) when you log in to the site.

To run the tests you can either set the environment variables in your `.bashrc`, `.profile` or `.bash_profile` and run:

`go test`

or on the command line:

`MAL_USERNAME=chris MAL_PASSWORD=foobar go test`

## TODO
- [ ] Documentation
- [ ] Command line interface

## Issues & Features
Of course, if you find a bug or need a feature added please feel free to submit an issue or pull request. I'm extremely eager to make this the best My Anime List API client in the universe (or something like that).
