include $(GOROOT)/src/Make.inc

TARG=atmos
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4
GOINSTALL=goinstall -dashboard=false

GOFILES=\
    atmos.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w atmos.go

install:
	${GOINSTALL} github.com/Arrow/GoGCM/atmos
  
