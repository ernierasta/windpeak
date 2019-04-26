# Windpeak
TES games (Skyrim, Oblivion, ...) author &amp; Nexus friendly modpack automation tool.

## Status
Total mess. Currently in experimentation phase. This repo is here only to allow me easy sync when switching computers/operating systems.

## Changelog

**2019-04-26**
Currently I have some drafts of GUI, Mod Organizer 2 support, Wrye Bash support. Nothing is working. My focus is now on codding modpack creator. Will be back when something will actually work.


## Planned features

- Mod Organizer 2 and Wrye Bash support for modpack creators and modpack users. It will be possible to create modpack with Wrye Bash and install it with Mod Organizer 2.

- Steps support, it will be possible to define steps by modpack author which allow for:
  - running some additional tools in the middle of installation, running game to check that everything is ok, etc.,
  - have alternatives to current step (for example Performance vs Quality option),
  - ability to skip step.

- Semi-automated modpack creating - everyone should be able to make modpack without spending week doing so.

- Reasonable upgrade modpack procedure.

## Why?

Yes, there are other tools to do that, I was inspired by Automaton(TODO: link), but I find it not as flexible as I would like. And it is written in C#.

To answer why make such tool. I have no time to play, usually I spend so much time installing mods for my games, that I have no time to play. Tools as this will cut installation time to minimum. But what is VERY IMPORTANT for me - mod authors or Nexus are NOT HURT in the process. User still have to download mod himself. Only install process is automated. In fact, this not different from what mod managers are doing.

