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

The next attempt is to parse a number. Anything in the range [a-f0-9] will parse as a number, these chaz being known as numbaz. Barring spaces, an Orc will eat numbaz until his stack is full, then push another one. A 16 bit Orc eats numbaz 4 at a type, a 32 bit Orc eats 8. 

Numbaz are always and only hexidecimal. No exceptions. 

If it is not a numba, it's a letta, and all lettaz are werdz. We spell everything in Orcish brutally, so that when we convert Forth words into Orcish werdz we don't break down in quivering confusion. 

Important note: an Orc will make numbaz until it can't: it switches to neutral as soon as this process fails, and attempts to keep chewing the input. So a stream like `1+` is identical to `1 +`, the Orc eats the 1, fails to make the `+` into a numba, and interprets the `+`. This is a boon to readability and the simplest possible implementation.

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

In a Core orc, we make grunt words manually. Our supervisor has the XT of DOCOL. available, and if we're exploring it's trivial to retrieve.

The header for an Orcish werd looks like ` [werd] [next] [head] | [code...] `. 	`next` is initially defined as EXIT. , which is rewritten with the address of `here` when a new werd is being defined. gruntz get interspersed, and will be forgotten if an Orc forgets any prior word. 

As a result ` ' fu 10 - `, for a 16 bit Orc will give DOCOL. if `fu` is a colon word, which it typically would be. 

Pseudocoding it up a bit ` h DOCOL. , ' D , ' * , EXIT. , ` would compile your basic square grunt. `h` is `here`, and since that's where DOCOL. ends up, we're good. 

Since Core Orcs have no assembler, direct words must be both constructed and linked to the bakpak 'manually'. This is of course a job for the supervisor environment. 

`;` merely compiles an exit, so to compile a real word, we'd say `: sq D * ; :`, such that `sq` would square the TOS as expected. 

`` ` `` causes the compiler to turn off for the next token. Entering compilation of a word has no effect on the stack, so `` ' D : fu ` , ; : `` will have the effect of manually compiling the XT of dup, making `fu` an inefficient synonym. 

There's no need immediate, postpone, or anything else that smacks of macros. Our stick shift is typically driven by robots, which are more than capable of macros. 

###Comments

`\` toggles the 'comment' mode. This actually writes to a special pad in memory in a circular queue, and the second `\` writes the range to a defined area of memory. The former is the `gab`, the latter is the `drp`, which is roughly the same as the Orc's state model. 

An umbilical system doesn't really need comments: this behavior plays a role in Orcish communication, both to sideband communication and to provide a dense block of data. This data can't have `\` in it, clearly, and should be composed of chaz, probably. 

####Aside

At first, we're going to fake a lot of this function just by jacking Retro. I plan to target Ngaro directly as soon as practical. I gather the Retro community would be stoked to see another language running on Ngaro, even if they think ours is weird. As indeed it is. 

###Stack Manipulation

`D` for **dup**, `s` for **swap**, `o` for **over**, `.` for **drop**.

###Math

`+,-,*,/,=,<,>` all behave as expected. In general, the extended glyphs such as `<>` should do the Right Thing, presuming an Orc understands them: if you think you know what a two-rune means, it's probably reserved.

Can we define `1+`? We cannot, this is a syntax error. We can and shall define `+1` with the same meaning, for some Orcs.

`A,O,N` are and, or, and xor. Them's your options.

###Control Flow




