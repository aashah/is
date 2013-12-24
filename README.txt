Installing:
Ensure you have golang installed onto your system. See {link here}

You can get the code for this program by using:
go get github.com/aashah/is

Afterwards, installation is made possible with:
go install github.com/aashah/is

The executable binary `is` can be found in your $GOPATH/bin

Goal:

Providing a user a clean way of working with modules. Since our primary focus
has been simplicity on the side of the user, abstracting away unecessary details
such as where modules are or how they get loaded, this command-line tool aims to
simplify dealing with module packages. A few features of this cli program are:

    get: The ability to get and build modules that reside in a version control
        system

    mv: Similar to get, but working with modules that pre-exist elsewhere on the
        file system

    run: Running the SDK

    build: Building a module

    chk: Checking a module for how "complete" it is

See `is help [subcommand]` for more information regarding these commands.

Author notes:

- A lot of this code is derived from the practices used by the golang team for
	implementing `go get` This includes their approach towards handling
	subcommands (such as get, build...etc) as well as their approach on creating
	and updating code retrieved from version control systems.
- Much of the code has been adapted to meet our needs, but the base structures
	are very close. More on the golang approach can be found at
	http://code.google.com/p/go


Roadmap:

    get: done
    mv: soon(tm)
    version: done
    run: soon(tm)
    build: done
    check: done