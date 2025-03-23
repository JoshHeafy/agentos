# COLORS
ccgreen=$(shell printf "\033[32m")
ccred=$(shell printf "\033[0;31m")
ccyellow=$(shell printf "\033[0;33m")
ccend=$(shell printf "\033[0m")

# Include other Makefiles
include ./scripts/make-files/build.mk

# SILENT MODE (avoid echoes)
.SILENT: all fmt test linter build