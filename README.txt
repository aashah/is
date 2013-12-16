Goal:

Providing a user a clean way of working with modules. Since our primary focus
has been simplicity on the side of the user, abstracting away unecessary details
such as where modules are or how they get loaded, a few features of this cli
program are:

    get: The ability to get and build modules that reside in a version control
        system

    mv: Similar to get, but working with modules that pre-exist elsewhere on the
        file system

    run: Running the SDK in a simple way

    build: Building a module

    chk: Checking a module for how "complete" it is

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