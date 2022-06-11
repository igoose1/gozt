SRC=main.go
BIN=bin
ZT-A=${BIN}/zt-a
ZT-L=${BIN}/zt-l
ZT-LL=${BIN}/zt-ll
ZT-N=${BIN}/zt-n
ZT-G=${BIN}/zt-g

.PHONY: all
all: ${ZT-A} ${ZT-L} ${ZT-LL} ${ZT-N} ${ZT-G}

${ZT-A}: ${SRC}
	go build -o $@

${ZT-L} ${ZT-LL} ${ZT-N} ${ZT-G}: ${ZT-A}
	cd ${BIN}; ln -s $(notdir ${ZT-A}) $(notdir $@)

.PHONY: clean
clean:
	rm -rf ${BIN}
