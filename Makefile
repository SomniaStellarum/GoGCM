include $(GOROOT)/src/Make.inc

TARG=gcm
DEPS=atmos
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
    gcm.go\
    
include $(GOROOT)/src/Make.cmd

format:
	${GOFMT} -w gcm.go
