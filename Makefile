NAME:=aoc2024
BIN_DIR:=./bin

PROG:=$(BIN_DIR)/$(NAME)

GOFILES:=$(shell find src/ -type f -name "*.go")

NOWDATE:=$(shell TZ="EST" date +%Y%m%d)
NOWDAY:=$(shell TZ="EST" date '+%e' | sed 's/^\s\+//')
ENDDATE:=20241225
DOCKERRUN=docker run --rm -i --env AOC_SESSION ${AOC_RUNOPTS} aoc2024:latest $(ELAPSEDOPTS)
ifdef ELAPSED
ELAPSEDOPTS:=-e
endif

.PHONY: build run run-all clean example build-run run-bare example-bare all today diy-run today-example

all: build

build: $(PROG)


$(PROG): $(GOFILES)
	DOCKER_BUILDKIT=1 docker build --target bin --output $(BIN_DIR)/ . 
	touch $(PROG)

build-run: $(PROG)
	docker build -f Dockerfile.run -t aoc2024:latest .

run: build-run $(PROG)
	@$(PROG) input $(DAY) | $(DOCKERRUN) $(DAY)

run-bare: $(PROG)
	@$(PROG) input $(DAY) | $(PROG) run ${AOC_RUNOPTS} $(ELAPSEDOPTS) $(DAY)

run-all: $(PROG)
	@if test "$(NOWDATE)" -lt "$(ENDDATE)"; then for day in `seq $(NOWDAY)`; do $(PROG) input $$day | $(DOCKERRUN) $$day; done; else for day in `seq 25`; do $(PROG) input $$day | $(DOCKERRUN) $$day;done;fi

today: build-run $(PROG)
	@$(PROG) input $(NOWDAY) | $(DOCKERRUN) $(NOWDAY)

today-example: build-run $(PROG)
	@cat examples/day$(NOWDAY).txt | $(DOCKERRUN) $(NOWDAY)

today-all: today-example today

clean:
	rm -f $(PROG)

example: $(PROG) build-run
	@cat examples/day$(DAY).txt | $(DOCKERRUN) $(DAY)

example-bare: $(PROG)
	@cat examples/day$(DAY).txt | $(PROG) run ${AOC_RUNOPTS} $(ELAPSEDOPTS) $(DAY)

diy-run: build-run $(PROG)
	@$(DOCKERRUN) $(DAY)

diy-run-bare: $(PROG)
	@$(PROG) run ${AOC_RUNOPTS} $(ELAPSEDOPTS) $(DAY)
