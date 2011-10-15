taller
=========

This is a golang templating lib. Here are it's assertions:

- I don't want to be limited by some idiotic ``people that write templates are stupid`` rule.
- I don't want to learn yet another template language.
- I want to use golang expressions if, for..
- I want all the power of indexing, attribute resolution and performing operations that golang offers.
- I want to type less stuff and do more.
- I want to include/extend other templates in my template. (think jinja)

Name
--------

**taller** (pronounced /ta.le/) is composed of **t** (template) and **aller** (the verb "go" in French).

Usage
-------

To know where the templates are **taller** requires a **TALLER_PATH** 
environment variable which is just like a normal unix **PATH** separated by ":". 

::

        package main

        import (
                "github.com/humanfromearth/taller"
                "os"
        )

        def main() {
                os.Setenv("TALLER_PATH", "/path/to/my/templates:/templates/")

                template := taller.TemplateFile("template.html")
                context := &taller.Context{"foo": "bar", "baz": 1}

                content := taller.Render(template, context)
                //taller is all bytes so we need to convert here to string
                println(string(content))
        }

For more examples of templates visit ``docs/examples/`` - TODO

How it works
-------------

Given a context and a template file translate the template directly into golang 
code executing it with the context and will returning the result.

Again ::

        template source -> golang code -> binary -> execute with context -> get the result

Why not just use a normal parser + reflect?
----------------------------------------------

I don't know how to write a decent parser/compiler so I'll use the golang 
compiler to do the job for me. Also it should be very fast. This of course 
makes the template language somewhat inflexible, but I think I can live with it.

[[Add benchmarks here]]

Parsing
-----------------

The rules are simple:

0. There are no spaces in the declaration of tags such as: [ if x == 1]. 
   Everything will begin with the actual tag: [if x[0] == 1]. 
   There are no 2 tags on the same line. These are the only style enforcing rules.

1. Find all [import ..], [extend ..] [include ..] in this order and perform the 
   required operations on them.

2. Find all blocks [block ...] and [endblock] and replace them in an order of 
   resolution the leafs have the highest priority.

3. Find all [if ...][endif], [for ...][endfor] and just replace them as 
   golang code. [if foo][endif] translates to if context["foo"] {} and so on.

4. Replace all \`expression\` with fmt.Fprintf(output, expression)

