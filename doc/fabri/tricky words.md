#Key Words for Inference in Forth

In keeping with the minimalism and grace of Forth, we want the type system to be as explicit as possible. When we can, we will build the type system in Fabri notation, to the side of the words as defined, just as a Forth environment is written (or at least bootstrapped) using as much Forth and as few primitives as practical. 

Some of the words are difficult to imagine a compact and general syntax for. This indicates to me that I don't fully understand the algorithm behind inferring through these words. They deserve disproportionate focus. 


## do ( xt -> mu )

`do` is more important in Retro than `execute` in Forths, being invoked in most control structures. 

The key to `do` is that we need to infer through to the type contract of the word being executed, when possible. Retro encourages quoting, which is straightforward from the type perspective, as it's simply an implicit word, albeit one you may nest within other quotes and within a definition. so `[ over swap dup ] do` will create a quote with the stack effect ( a b -> a a b b ), while `do` will actually have the stack effect `( a b xt -> a a b b )`. To be clear, the stack effect of `]` is `( nil -> xt )`.

I'm still brushing up on how Retro actually builds quotes in Ngaro. Regardless, I presume there's a reliable way to tell an XT from another address, meaning we can perform the usual interpret or compile time check against the word contract. 

## if ( flag xt xt -> mu | mu )

`if`, in Retro, executes NOS if the flag is -1 and TOS if the flag is 0. 

This brings up iterators, one place where Completeness suggests we cannot infer all cases. We can confirm the case where the quote being iterated has no net effect on the stack, which is usually what we want to do: looping a variable amount of junk onto the stack, then looping it off again, is poor form in Forth. 

The more complex combinators are hopefully straighforward applications of these two, though 'if' is likely not the actual primitive that needs fanciness built onto it. 

## dup ( a.cell -> a.cell a.cell )

Duplicate tracking will prove important. Any assertion made on one cell of a pair applies automatically to the other. Very occasionally, one might not want this; normally, it's the right thing.

This intersects with memory:

## memory: @ and ! ( cell adr -> !cell ) & ( adr -> @cell )

Inferring ordinary stack effects is conceptually straightforward. What we're able to do with memory will define the limits of Fabri. One of several reasons I've chosen Ngaro as the platform: the memory model is very easy to understand.

Let's start with the idea of a variable, which is initially of type `@cell`, meaning roughly "the use of `@` will retrieve a cell from this variable". Well and good, but we can call `allot` and `create` primitively. In most cases, having a useful type system around should encourage defined ranges. 

We still have to handle the case where a range of memory is of type `@mu`, which is not a very good annotation but I haven't determined a vocabulary for range yet. For the most part, we'll want to just give up.

However, we can use dup tracking in many cases to carry the range over. Storing data to generic memory then losing all reference to it is a memory leak by definition. Normally we'll have a copy of the address stored somewhere. If the address is named, it's a variable, and we track those explicitly. If it's not, our dup tracker will change the type of the duplicate address from `@mu` to `@type`, so we're good if we try to call a word restricted to `type` on that address. 
Inevitably, we need to encourage certain stylistic choices that support inference. It will be perfectly possible to write programs that are perverse and cannot be inferred. This is a code smell, not a reason to introduce dependant typing in the rabbit-hole sense. 





