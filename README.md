...TODO the rest of the readme...

### Example ALSA Config
create a virtual ALSA device that multiplexes to the soundcard and a virtual loopback device
```
pcm.!default {
    type hw
    card sndrpihifiberry
}

pcm.!default {
  type asym
  playback.pcm {
    type plug
    #slave.pcm "output"
    slave.pcm "loopdac"
  }
  capture.pcm {
    type plug
    slave.pcm "loopin"
  }
}

ctl.!default {
  type hw
  card sndrpihifiberry
}


# TODO: is this necessary
pcm.front-plug {
    type plug
    slave.pcm "front:CARD=sndrpihifiberry,DEV=0"
}


# Setup a multiplex device which sends to both the hifiberry and the loopback device

# a virtual device that copies stereo data into 4 channels for the mux below
pcm.loopdac {
  type route
  slave.pcm loopdacmux
  slave.channels 4

  ttable {
    0.0= 1
    0.2= 1
    1.1= 1
    1.3= 1
  }
}

# a virtual device with 4 channels, two to the DAC
# and two to loopback
pcm.loopdacmux {
  type multi
  slaves {
    a {
      pcm "plughw:sndrpihifiberry,0"
      channels 2
    }
    b {
      pcm "hw:Loopback,0,0"
      channels 2
    }
  }
  bindings.0.slave a
  bindings.0.channel 0
  bindings.2.slave b
  bindings.2.channel 0
  bindings.1.slave a
  bindings.1.channel 1
  bindings.3.slave b
  bindings.3.channel 1
}


pcm.loopin {
 type dsnoop
 #ipc_key 1111
 #ipc_key_add_uid false
 slave.pcm "hw:Loopback,1,0"
}
```
