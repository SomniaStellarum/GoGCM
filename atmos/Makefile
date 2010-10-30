include $(GOROOT)/src/Make.inc

TARG=atmos
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
    atmos.go\

include $(GOROOT)/src/Make.pkg

format:
${GOFMT} -w atmos.go

