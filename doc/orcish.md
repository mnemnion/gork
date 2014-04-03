#Orcish

This is a quick guide to Orcish, intended to aid the initial implementation on the Ngaro virtual machine.

The detailed discussion is [elsewhere](https://github.com/mnemnion/imp/tree/master/Forge/markdown/ArdForth). 

##Fundamentals

Orcs are envisioned as a tribe of microcontroller systems divided by a common language. 

That language is Orcish. It is Forth boiled down. 

##Language

Orcs do several things backwards from the Forth we're used to. This is a fertile place to start.

Orcs hear one byte at a time. The first thing an Orc will do with a byte is try and ignore it. If it's not a printable ASCII character, a space, or (maybe) a newline, an Orc is casually oblivious. This may cause confusion, discussed later.

Any byte which succeeds in this first pass is a cha. Orcs hear and speak in cha. 

The next attempt is to parse a number. Anything in the range [a-f0-9] will parse as a number, these chaz being known as numbaz. Barring any other character, an Orc will eat numbaz until his stack is full, then push another one. A 16 bit Orc eats numbaz 4 at a type, a 32 bit Orc eats 8. 

Numbaz are always and only hexidecimal. No exceptions. 

If it is not a numba, it's a letta, and all lettaz are werdz. We spell everything in Orcish brutally, so that when we convert Forth words into Orcish werdz we don't break down in quivering confusion. 

Important note: an Orc will make numbaz until it can't: it switches to neutral as soon as this process fails, and attempts to keep chewing the input. So a stream like `1+` is identical to `1 +`, the Orc eats the 1, fails to make the `+` into a numba, and interprets the `+`. This is a boon to readability and the simplest possible implementation.

There are gotchas! `1D` is a one and a dup, not a numba. 

Note further: When an Orc chews `1e7`, it chews `1` first, puts it on the stack, chews `e`, putting `1e` on the stack through simple multplication, and so on. When the stack is full, it pushes a fresh number: `ffffaaaa` gives ( ffff aaaa -- ) on a 16 bit system. 

##Werdz

All Orcs are familiar with one letta werdz. The meaning is the same, insofar as possible across architectures. 

The one letta werdz are the API, the DNA of Orcish kind. The minimum necessary for computation, stack manipulation, word definition, and the various memory accesses. 

This is a different set of words from the eForth core or other more literal groups. Anything which we don't care to expose in the core API is a grunt, similar to a quote or headless word. An Orcish monitor can disassemble these words and provide access to them, if one is programming through the Forth interface (and one normally will). 

###Two Letta Werdz

A one letta werd must be followed by a space, as must (most) two letta werdz. 

The tawka (the Orcish interpreter) does a certain amount of heavy lifting. In this case, it detects any of `#`,`,` and `%` as one of the two lettaz, in either order. These are automatic `#constants`, `,variables` and `%values`, respectively. An Orc has 155 of each. 

werdz consisting of exactly two letters (not lettaz, letters) are reserved for user words. They may not begin with a numba, or they will be parsed as one, and possibly fail. There are 1,536 of these available. 

The remaining 5,331 (give or take) werdz are the Orcish library. Some effort is made to contain the chaos.

Orcs may be taught two letta werdz, and it's likely even the most primitive Core Orcs will be expected to know a few. The structure (the bakpak) which contains the two letta werdz must exist, and the core API has `:` to add werdz to it. 

Again, this is unlike most Forth, in that the dictionary is a forward-linked list and an Orc will refuse to relearn a word. It can be induced to forget one, which causes it to forget everything after. 

###Compilation model

The compiler is a stick shift. I'm still exploring precisely what this means; there's a lot of good MachineForth spinoffs to look at here.

`:` is the compiler toggle. In interpretation mode, it reads the next token. If it can define it, compiler mode is retained, otherwise the Orc complains audibly (which is unusual) and turns back to interpreter. 

If `:` is encountered in compile mode, it turns back to interpreter. Both `;`, which compiles an exit, and `:` must be present to end a definition. 

Grunts, `:noname` in classic Forth, turn on the compiler with `|`. This drops a DOCOL and proceeds as usual. If you want the XT on the stack, use `h` for **here** before you get started. Grunts terminate with `; |` though in compilation mode `:` and `|` have identical effects and both work. ` h | D * ; | ` puts an address on the stack and makes your classic square-shaped grunt at that location. ` 3 h | D * ; | x `, for **execute**, leaves 9 on the stack. 

Since Core Orcs have no assembler, direct words must be both constructed and linked to the bakpak 'manually'. This is of course a job for the supervisor environment. 

`;` merely compiles an exit, so to compile a real word, we'd say `: sq D * ; :`, such that `sq` would square the TOS as expected. If we leave off the second `:` we'll keep compiling unreachable words.

`` ` `` causes the compiler to turn off for the next token. Entering compilation of a word has no effect on the stack, so `` ' D : fu ` , ; : `` will have the effect of manually compiling the XT of dup, making `fu` an inefficient synonym. 

There's no need immediate, postpone, or anything else that smacks of macros. Our stick shift is typically driven by robots, which are more than capable of macros. 

###Comments

`\` toggles the 'comment' mode. This actually writes to a special pad in memory in a circular queue, and the second `\` writes the range to a defined area of memory. The former is the `gab`, the latter is the `drp`, which is roughly the same as the Orc's state model. 

An umbilical system doesn't really need comments: this behavior plays a role in Orcish communication, both to sideband communication and to provide a dense block of data. This data can't have `\` in it, clearly, and should be composed of chaz, probably. 

####Aside

At first, we're going to fake a lot of this function just by jacking Retro. I plan to target Ngaro directly as soon as practical. I gather the Retro community would be stoked to see another language running on Ngaro, even if they think ours is weird. As indeed it is. 

A simple example: instead of implementing a forward-linked bakpak, we'll just look in the Orc's 'chain' for the word, and refuse to define it if already present. We end up with the same behavior. 

###Stack Manipulation

`D` for **dup**, `s` for **swap**, `o` for **over**, `.` for **drop**, `r` for **rot**.  `{` and `}` work like `r>` and `>r`, pushing and popping across the `sak` and `bag`. Those are our Orc's data stack and return stack, naturally.

###Math

`+,-,*,/,=,<,>` all behave as expected. In general, the extended glyphs such as `<>` should do the Right Thing, presuming an Orc understands them: if you think you know what a two-rune means, it's probably reserved.

Can we define `1+`? We cannot, nor need we. It's parsed as `1+`, as is `0=` the same as `0 =`. 

I think we probably want `~` to mean modulus. 

`A,O,N` are and, or, and xor. Them's your options. We do need unary not, defined correctly, which is probably `^`.

###Control Flow

I'm wondering how stripped down we can get this.

`i` for **if**, `E` for **else**, `t` for **then**, which is optional if both code blocks exit. I do want to support nested conditionals even within the core, for a few reasons. 

I do propose an unusual counted loop. `l` takes ( XT index -- ) and executes from 0 to index-1 times, storing the index in an implementation-specific register which may be manually read. So ` 3 ' D 10 l` would leave 11 `3`s on the stack. 

I gather this will be small and fast: a custom DOCOL. that skips the header, executes the word, and iterates on EXIT by incrementing the return pointer and jumping back to the start. 

Conditional loops benefit from our minimalist language. Anything between `(` and `)` is repeated so long as the TOS is true, `( 1 )` will hang your machine nicely. ` 1 ( D 1+ D 20 = ) ` will leave a series of numbers up to 20 on your stack, and maybe blow it; Orcs shouldn't be counted on to be deep in the `sak` department.

Do we need more control structures? 
These nest, of course. It's a nice idiom, actually. 

###Memory 

Orcs handle any token containing `#`, `%`, or `!` in a special fashion. A token such as `Q#`, when first encountered, has the stack effect ( value -> nil ); subsequent uses of a# will have the effect ( nil -> value ), such that these behave somewhat like constants. 

Orcs are protective of the pak. It lives in Flash, and it's brittle as a result. We don't provide access to the Flash words directly. `:`, `|`, `,`, and `#` all write to the Flash, and most anything reads from it. 

`%` words access the spleen, our EEPROM in Harvard Yard. `%` itself reads eeprom ( adr -> value), while `Q`, a nonce word, performs a write ( value adr -> nil ).

Words like `Q%` or `%Q` are like values, set at definition so ` 23 %t ` stores 23 somewhere useful in the EEPROM. The default behavior of `%t` is to read EEPROM. To write to it, `24 ' %t q`, q for nothing in particular. 

Yes, we need two whole words here, because values are stored in the pak as "read from this area of eeprom". So the XT itself doesn't point to the EEPROM, though the address is near it, in an implementation dependent way. 

`!` words are the same but in memory, and calling them places the memory address on the stack: `23 m!` makes it and `m! @` reads a 23 onto the sak. `!` itself and `@` write and read as usual: AVRS can read the same way across both memories, though eeprom is separate and writing is very much so. `!` on any address in the pak shouldn't work, while ` ' fu @ ` would read the XT for DOCOL. from the header of `fu`. 

###Communication

Orcs expect to be in the company of other Orcs. They also glumly accept being bossed around by users, because they can't tell the difference. 

The most impressive capacity Orcs have is self-report. Let's start with the basics.

When you tell an Orc to do something, and he understands you, he does it. First, he tells you he understood you, by repeating the command back as a comment. So the user says `10`, the Orc says `\ 10 \ ` back. The user says `D`, the Orc says `\ D \ ` then dups the `10`. 

If we have, instead of a user, another Orc, the commanding Orc will treat the comments as comments. 

If an Orc doesn't understand you, it says nothing. You may imagine it glaring menacingly. It also increments the confusion byte. This goes down when you make sense again, until the Orc is happy. 

Orcs will happily read their memory out to you, using `r ( start end -> nil "range" )`. This is hex, wrapped in a comment and pretty printed with spaces and newlines. 

The real treasure is `?`, which prompts the Orc to search for the next token. `? fu` will earn you a comment containing the full Orcish definition of fu. This is a superpower for such a small beast. 

Lets say we have `: fu 34 + 12 / ; :` as a random definition. `? fu` will compell the Orc to utter the fell words `\ : fu 34 + 12 / ; : \ ` just as pretty as you please. Edge case: if you've stuck the actual function `\` in something, and `` ` ' \ `` will in fact do this, your comment will end in a bad place. So don't do that, or account for it. There's a newline after the concluding comment, and `?` never generates a newline otherwise, so there's that: this kind of nonsense is unlikely.

If the word is direct, the Orc provides a hex dump, word aligned, up to the NEXT. If the next token is a numba, the Orc attempts to read a grunt from that address: if it's a DOCOL, a disassembly, direct gets a dump, otherwise, it glares at you. So `? ab23` tries to read the grunt from ab23. 

`h` 'hears' a single byte, putting it on the stack, while `S` says whatever's on top, popping it: only the low byte is sent. `Z` interprets the stack as hex and says it accordingly. `w` polls for a byte, hearing it if there's input; `h` blocks until it hears something. `h` and `S` do not filter the byte into a cha. 

The eye and ear port values are stored in the spleen; changing them is perforce implementation dependent. There will be an i/o library for fast multiport programming, likely centered around tokens with `\` in them. 

###Integrity

`C`, stack effect `( adr count -> checksum )`, does an adler32 over the region of memory. 

I have some notions of Orcs using multiplexing to get around bad environments and to get away from bad input, but those are somewhat half-baked. 