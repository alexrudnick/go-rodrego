go-rodrego
==========

This is a reimplementation of the RodRego register machine in Go, and hopefully
performs the same computations as the version from Tufts, which is available
[for download](http://sites.tufts.edu/rodrego/) and [on
GitHub](https://github.com/TuftsUniversity/RodRego) as well.

To understand register machines and what's going on generally, take a look at
Dan Dennett's
[The Secrets of Computer Power
Revealed](http://sites.tufts.edu/rodrego/files/2011/03/Secrets-of-Computer-Power-Revealed-2008.pdf),
or his very nice book "Intuition Pumps and Other Tools for Thinking".

I made go-rodrego mostly just as a small exercise for learning Go, but also
because the Tufts version seems to be in RealBASIC (now called
[Xojo](http://en.wikipedia.org/wiki/Xojo) apparently) and I couldn't figure out
how to make it go on Linux without downloading... uh... Xojo.

Anyway! This is a tiny interpreter for the RodRego register machine, and I hope
somebody will find it interesting.

## installation
You'll need Go installed already -- I have Go release 1.2.1 on Ubuntu, but later
versions will probably work.

Get the package like so:

    $ git clone https://github.com/alexrudnick/go-rodrego.git

Then build:

    $ cd go-rodrego
    $ go build rodrego.go

And there you go! You can now run ./rodrego

## running RodRego programs
You can specify both the program to run and the initial state of the registers
with names of files on the command line.

 * The ``-program`` argument (obligatory) is for specifying the RodRego program to run. Give it a filename.
 * The ``-values`` argument (optional) is for specifying the initial state of the registers for the register machine.
 * The ``-step`` flag makes the interpreter step through a program one instruction at a time.

For example:

    $ ./rodrego -program testprograms/add.rgo -values values.txt

This will run the ``add.rgo`` file with initial register values set by
``values.txt``. If your ``values.txt`` looks like this:

    2 10
    3 5

... (which sets register 2 to have value 10, and register 3 to have value 5),
then the end result should be their sum, 15, in the output register.

## file formats
Interestingly, RodRego ships with sample programs that have old-style (pre-Mac
OS X) Mac line endings; just a carriage return character, rather than the more
familiar Unix/Mac line endings (line feed) or DOS/Windows-style CRLF.

RodRego supports programs with any of these line endings, which is slightly
interesting because the Go standard library routines don't expect by default to
see old-style Mac line endings.

### rodrego programs
To learn how these programs work, take a look at Dan Dennett's excellent
[The Secrets of Computer Power
Revealed](http://sites.tufts.edu/rodrego/files/2011/03/Secrets-of-Computer-Power-Revealed-2008.pdf).
He does a fabulous job of explaining what's up here.

RodRego programs have one instruction per line, of the form LINENUMBER
INSTRUCTION ARGUMENTS. Instructions can be one of ``INC``, ``DEB`` or ``END``
(any case is fine).

  * INC takes two arguments, a register number to increment and a line to jump to next.
  * DEB is "decrement or branch", and it takes three arguments: the register number to degrement, a line number to jump to if that can be done, and a line number to jump to if the register is already at 0. This is the only way to do a conditional branch in the language.
  * END ends the program.

Line numbers are canonically natural numbers in base-10, but in this
implementation you can actually use any kind of label as long as it doesn't have
whitespace in it. You could use ``aubergine`` as a line number if you want.

Also comments are a thing; a line where the first non-whitespace character is a
``#`` is treated as a comment.

### register values
One register can be set per line; put the register number and then its value.

Both must be base-10 nonnegative integers.

----

## OK, that's it! Have fun!
