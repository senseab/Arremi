#ifndef ARREMI_ALSA_WRAPPER_H
#define ARREMI_ALSA_WRAPPER_H

#include <alsa/asoundlib.h>

static snd_seq_t *seq_handle;
int seq_client, seq_port;

int new_client(char*);
int new_port(char*);
int send_data(char*, int);
#endif