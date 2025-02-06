# go-bouncingball
Bouncing Ball demo in GO

One of my favourite things to write in any new language I come across is program to bounce a ball around the screen.

Please note with this repo, I am EXTREMELY NEW to GO.  While this demo does work, and seems to work well, it's written by a GO novice, and as and when my GO skills improve i'll refactor this.

This is a console mode application.  For controlling console i/o I use https://github.com/rivo/tview/wiki 

The bouncing of the ball should auto-correct itself for any change in terminal window size

![Sample Output](./go-bouncingball.gif)

Controls are shown to:
* [ESC] exit the demo application
* [+] to add a ball (maximum of 10)
* [-] to remove a ball (down to 1)
* [CursorUp] Increase the speed of all balls
* [CursorDown] Decrease the speed of all balls

Each ball is made up of a ball 'head' and a 'tail'.  Each new ball is started in a random position, a random direction and a random colour.

I've chosen characters to represent a ball from a Unicode character set, with the characters of the tail getting smaller.  The colour of the tail fades too, which helps to portray a streaking effect as it would been viewed on an old phosphor CRT screen back in the day.

The TView library that is used for the console i/o is cross platform, so this should work just as well on Windows and Linux as it does on my Mac

