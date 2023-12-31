#define _GNU_SOURCE
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mount.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

void check (int res, const char *msg) {
  if (res) {
    perror(msg);
    exit(1);
  }
}

int main (int argc, char **argv) {
  if (argc < 2) {
    fprintf(stderr, "Usage: %s [ns name] [cmd args...]\n", argv[0]);
    exit(2);
  }
  const char *base = "/tmp/netns/ns/";
  char *devname = argv[1];
  char buf[65];
  if (strlen(devname) > sizeof(buf) - strlen(base) - 1) {
    fputs("Device name too long\n", stderr);
    exit(2);
  }
  strcpy(buf, base);
  strcat(buf, devname);

  int fd = open(buf, O_RDONLY);
  check(fd < 0, "open");
  check(setns(fd, CLONE_NEWNET), "setns");

  if (argc == 2) {
    const char *sh = getenv("SHELL");
    if (!sh) sh = "/bin/bash";
    check(execl(sh, sh, "-l", NULL), "execl");
  }

  argv += 2;
  check(execvp(argv[0], argv), "execvp");
}
