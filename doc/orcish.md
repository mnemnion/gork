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

##Werdz

All Orcs are familiar with one letta werdz. The meaning is the same, insofar as possible across architectures. 

The one letta werdz are the API, the DNA of Orcish kind. The minimum necessary for computation, stack manipulation, word definition, and the various memory accesses. 

This is a different set of words from the eForth core or other more literal groups. Anything which we don't care to expose in the core API is a grunt, similar to a quote or headless word. An Orcish monitor can disassemble these words and provide access to them, if one is programming through the Forth interface (and one normally will). 

###Two Letta Werdz

A one letta werd must be followed by a space, as must (most) two letta werdz. 

The tawka (the Orcish interpreter) does a certain amount of heavy lifting. In this case, it detects any of `#`,`,` and `%` as one of the two lettaz, in either order. These are automatic `#constants`, `,variables` and `%values`, respectively. An Orc has 155 of each. 

werdz consisting of exactly two letters (not lettaz, letters) are reserved for user words. They may not begin with a numba, or they will be parsed as one, and possibly fail. There are 1,536 of these available. 

The remaining 5,331 (give or take) werdz are the Orcish library. Some effort is made to contain the chaos.

