#include "stdio.h"
#include <sys/socket.h>
#include <dlfcn.h>
#include <stdlib.h>

#define SETUP_SYM(X) do { if (! true_ ## X ) true_ ## X = load_sym( # X, X ); } while(0)

static void* load_sym(char* symname, void* proxyfunc) {
	void *funcptr = dlsym(RTLD_NEXT, symname);

	if(!funcptr) {
		fprintf(stderr, "Cannot load symbol '%s' %s\n", symname, dlerror());
		exit(1);
	} else {
		fprintf(stderr, "loaded symbol '%s'" " real addr %p  wrapped addr %p\n", symname, funcptr, proxyfunc);
	}
	if(funcptr == proxyfunc) {
		fprintf(stdout,"circular reference detected, aborting!\n");
		abort();
	}
	return funcptr;
}
typedef ssize_t (*sendto_t) (int sockfd, const void *buf, size_t len, int flags,
			     const struct sockaddr *dest_addr, socklen_t addrlen);

sendto_t true_sendto;

ssize_t sendto(int sockfd, const void *buf, size_t len, int flags,
	       const struct sockaddr *dest_addr, socklen_t addrlen) {
    fprintf(stdout,"hello world\n");
    return 0 ;
}


static void setup_hooks(void) {
	SETUP_SYM(sendto);
}

