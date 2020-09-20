target=libhooks.so


all: ${target}

.PHONY: clean

clean:
	rm -rf *.o *.so

${target}: hooks.c
	gcc -shared -o ${target} -fPIC -D_GNU_SOURCE hooks.c

