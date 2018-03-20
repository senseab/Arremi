#include <alsa/asoundlib.h>
#include "alsa_wrapper.h"

int new_client(char* client_name) {
    int status;
    status = snd_seq_open(&seq_handle, "hw", SND_SEQ_OPEN_OUTPUT, 0);
    if (status < 0) {
        return 1<<16|status;
    }

    seq_client = snd_seq_client_id(seq_handle);
    if (seq_client < 0) {
        return 2<<16|status;
    }

    status = snd_seq_set_client_name(seq_handle, client_name);
    if (status < 0) {
        return 3<<16|status;
    }

    // Success
    return 0;
}

int new_port(char* port_name) {
    seq_port = snd_seq_create_simple_port(seq_handle, port_name, 
                SND_SEQ_PORT_CAP_READ|SND_SEQ_PORT_CAP_SUBS_READ,
                SND_SEQ_PORT_TYPE_MIDI_GENERIC|SND_SEQ_PORT_TYPE_APPLICATION);
    if (seq_port < 0) {
        return seq_port;
    }

    // Success
    return 0;
}

int send_data(char* data, int length) {
    int status;
    snd_midi_event_t* midiParser;

    status = snd_midi_event_new(length, &midiParser);
    if (status < 0) {
        return 1<<16|status;
    }

    snd_seq_event_t* ev;
    snd_seq_ev_clear(ev);

    status = snd_midi_event_encode(midiParser, data, length, ev);
    if (status < 0) {
        return 2<<16|status;
    }

    snd_seq_ev_set_direct(ev);
    snd_seq_ev_set_source(ev, seq_port);
    snd_seq_ev_set_dest(ev, SND_SEQ_ADDRESS_SUBSCRIBERS, 0);    

    status = snd_seq_event_output_direct(seq_handle, ev);
    if (status < 0) {
        return 3<<16|status;
    }

    return 0;
}