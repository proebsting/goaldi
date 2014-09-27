#!/bin/sh
#
#  goaldi file [ arg... ] -- compile and execute (interpret) Goaldi program
#
#  Assumes that gtran and gexec are in the search path.

I=${1?"usage: $0 file [ arg... ]"}
shift

export COEXPSIZE=300000
exec gtran preproc $I : yylex : parse : ast2ir : optim -O : json_File : stdout |
  gexec "$@"
