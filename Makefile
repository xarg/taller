include $(GOROOT)/src/Make.inc
TARG=github.com/humanfromearth/taller
GOFILES=\
	compiler.go\
	parser.go\
	utils.go\
	taller.go
include $(GOROOT)/src/Make.pkg
