include $(GOROOT)/src/Make.inc

TARG=timestep
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4
GOINSTALL=goinstall -dashboard=false 

GOFILES=\
    timestep.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w timestep.go
	${GOFMT} -w timestep_test.go

install:
	${GOINSTALL} github.com/Arrow/GoGCM/util/timestep
    
