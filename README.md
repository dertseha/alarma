# Alarma - A bedside clock with music capabilities

This is a personal project. I wanted to have a bedside clock with alarm functionality.

Primary requirements:
* Full-screen application, emitting as few light as possible
* Show the current time in hours and minutes
* Within configured time-spans, random audio files from a directory shall be played
* Simple user-interaction to switch "alarm" functionality on/off

Secondary requirements:
* Single time-spans shall have capability to be disabled. Sometimes a different time is required without needing much config work.


## Target hardware

This application shall run on a touch-capable display device. Single-Touch input shall be interpreted as mouse events.
My initial target is a "[NanoPi 2 Fire](http://wiki.friendlyarm.com/wiki/index.php/NanoPi_2_Fire)", combined with an LCD.

## Application

Each time-span has a separate folder configured. This way I can configure one time-span with music, and a second one with some classic alarm sound.

## Dependencies

* sox audio player (```play``` command) -- hardcoded for now.
* OpenGL 2.1 capabilities


## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
