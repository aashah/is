COMMAND LINE STRUCTURE:
Prefix: is
Subcommands: {
    get: [-v verbose] [-b build] [modules]
    mv: [-v verbose] [-b build] [directories]
    run: [-v verbose]
}

- A lot of this code is derived from the practices used by the golang team for
	implementing `go get` This includes their approach towards handling
	subcommands (such as get, build...etc) as well as their approach on creating
	and updating code retrieved from version control systems.
- Much of the code has been adapted to meet our needs, but the base structures
	are very close. More on the golang approach can be found at
	http://code.google.com/p/go


Roadmap:

    get:
        retrieve module: done
        check integrity: soon(tm)
        build: soon(tm)
    mv: soon(tm)
    version: done
    run: soon(tm)
    build: soon(tm)
    check: soon(tm)