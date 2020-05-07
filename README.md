# tf-auto-document
A tool which automatically builds documentation for Terraform modules.

And before you say it, I know there are other tools that do this.  I created this because:

* I wanted to learn about Go
* I wanted the documentation to include git tag (release) information

## Example
I use this tool on my Terraform modules hosted here: https://github.com/richardjkendall/tf-modules

## Module structures
This tool expects your repository to be organised as follows

```
root
 |-modules
   |-module 1
     |-main.tf
     |-variables.tf
     |_...
   |-module 2
     |-...
```

It will scan each module and find the variables and outputs and include those in the documentation.

It will look in `main.tf` for a comment at the start of the form

```
/*
title: example
desc: This is an example module
partners: another-module-I-work-with
depends: a-module-I-depend-on
*/
```

The title must contain only lower/uppercase characters A-Z or hyphens.

## How to use

### Build yourself

1. Clone this repo
2. Build with `go build`
3. Run tool pointing it at a folder containing a git repository containing Terraform modules.  `./tf-auto-document ../tf-modules`

### Use a pre-built binary

1. Download a a binary from the Releases tab (select one for your OS)
2. Decompress it
3. Run tool pointing it at a folder containing a git repository containing Terraform modules.  `./tf-auto-document ../tf-modules`
