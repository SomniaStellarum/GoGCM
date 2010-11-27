include $(GOROOT)/src/Make.inc

TARG=solar
DEPS=../util/timestep
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4
GOINSTALL=goinstall -dashboard=false

GOFILES=\
    solar.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w solar.go
	${GOFMT} -w solar_test.go
	
install:
	${GOINSTALL} github.com/Arrow/GoGCM/solar

