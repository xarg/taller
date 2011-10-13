taller
=========

A golang templating lib. Here are it's assertions:

- I'm not stupid 
- I don't want to be limited by some idiotic ``people that write templates are stupid`` rule.
- I don't want to learn yet another template language.
- I want to use golang expressions if, for..
- I want all the power of indexing, attribute resolution and performing operations that golang offers.
- I want to type less stuff and do more.
- I want to include other templates in my template. (jinja?)
- I want to extend templates.

Implementation
-----------------

The rules are simple:

0. There are no spaces in the declaration of tags such as: [ if x == 1]. 
   Everything will begin with the actual tag: [if x[0] == 1]. 
   There are no 2 tags on the same line. These are the only style enforcing rules.

1. Find all [import ..], [extend ..] [include ..] and perform the 
   required operations on them.

2. Find all blocks [block ...] and [endblock] and replace them with actual 
   golang function calls with the same names in the main loop. 
   Create the functions.

3. Find all [if ...][endif], [for ...][endfor] and just replace them as 
   golang code. [if foo][endif] translates to if context["foo"] {} and so on.

4. Replace all \`expression\` with fmt.Fprintf(output, expression)?
