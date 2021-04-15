# CordCGI
Do one thing, and one thing terrible.

![Image of Trollface inventing CordCGI](./cordcgi.jpg)

---

I took the worst parts of WebCGI and turned it into a web framework:

Make a directory `cgi-bin`, add a bunch of programs to that.
At the moment, the only meaningful thing you can do with this is reply to messages.
But that's OK!
Write your hello world in bash:

```
$ cat > cgi-bin/hello <<EOF
#!/bin/sh
echo "<@${DISCORD_MESSAGE_AUTHOR}> Hello, World!"
EOF
$ chmod +X cgi-bin/hello
```

Or maybe you like Python:

```
#!/usr/bin/env python3
import os

user = os.environ['DISCORD_MESSAGE_AUTHOR']
print(f"<@{user}> Hello, World!")
```

Perl? Why not!

```
#!/usr/bin/perl
use Env qw(DISCORD_MESSAGE_AUTHOR);

print "<\@$DISCORD_MESSAGE_AUTHOR> Hello, World!";
```

Go crazy. Go stupid!!!

```
#!/usr/bin/tcc -run
#include<stdio.h>
#include<stdlib.h>

int
main(void)
{
	printf("<@%s> Hello, World!\n", getenv("DISCORD_MESSAGE_AUTHOR"));
	return 0;
}
```

This is really a bad way to write software.

## Setup

Copy `settings.json.example` to `settings.json`.
Write your client ID, bot token, and preferred command prefix.
Set `CgiBin` to the path to the directory that calls your code.
I recommend setting it to `./example-cgi-bin` and trying out all the commands there.
Relative paths are relative to the working directory when you run the program.

`go get github.com/bwmarrin/discordgo`, `go run .` and *bam!* Discord bot!

Your `CgiBin` directory (which I'm assuming is `cgi-bin/` for now) would be a collection of programs -- scripts, binary executables, anything -- that would be run for each command.
If someone, for example, says `!add 10 20` (assuming the prefix is `!`), and there is an executable in `cgi-bin` called `add`, it is run with the argv arguments `10` and `20`.
One interesting thing to note is that subdirectories kind-of work as subcommands.
The command `!directory/command` *is* a valid pathname, and if `cgi/directory/command` exists, that would be executed.

When a program is run, it will be provided the entire message's content into stdin, it is provided message data as environment varaibles, and it will take all stdout as a reply.
The environment varables are:

- `DISCORD_MESSAGE` - the message's ID
- `DISCORD_MESSAGE_AUTHOR` - the message's author's ID
- `DISCORD_MESSAGE_AUTHOR_AVATAR` - the author's avatar ID
- `DISCORD_MESSAGE_AUTHOR_LOCALE` - the author's preferred locale
- `DISCORD_MESSAGE_AUTHOR_USERNAME` - the author's username
- `DISCORD_MESSAGE_CHANNEL` - the message's channel's ID
- `DISCORD_MESSAGE_GUILD` - the message's guild's ID
- `DISCORD_MESSAGE_UNIXNANOS` - the timestamp of the message, as the nanoseconds since Unix epoch.

## Reflection

In retrospect, if the only discord things you need to do is read the message calling commands and reply to them, this would probably be a *really* good way to prototype it.
Once your bot needs to become any more rich, though, this probably falls apart.
Overall, I recommend *against* using this for serious bot development unless you're doing it for laughs.
