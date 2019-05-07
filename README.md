# Windpeak
TES games (Skyrim, Oblivion, ...) author &amp; Nexus friendly modpack automation tool.

## Status
Total mess. Currently in experimentation phase. This repo is here only to allow me easy sync when switching computers/operating systems.

## Changelog

**2019-04-26**
Currently I have some drafts of GUI, Mod Organizer 2 support, Wrye Bash support. Nothing is working. My focus is now on codding modpack creator. Will be back when something will actually work.
**2019-04-30**
Wasted some time by investigating EOF error from Nexus websocket - turns out, that reusing UUID will make nexus just drop connection. It was quite random (probably there is load balancer). I have first production code in **nexusapi** dir. I had implemented NXM handler on Windows and Linux, but is not here yet, more work needed. Now I need to implement automatic mod definitions(from MO2 and WB) and start making modpack creator GUI.
**2019-05-07**
Creator GUI is almost complete for now. It has almost all functions to test real modpack creation. I had needed to workaround few libui/ui library missing features. I could use different UI lib, but I like this one and mainly - it is multiplatform. GUI code is nothing I am proud about, but for now it is good enough. We will see after user GUI will be written. So now, I will probably design and write modpack configuration files parser boiler-plate code.
BTW: I was also looking into the Morrowind modding and while it has some quirks, it should be possible to support also Morrowind and especially OpemMW, where Linux support would be useful. I should collect some notes on this topic.

## Planned features

- Multiplatform - Linux and Windows will work. Darwin will remain untested/unsupported until someone will want to test it there.

- Mod Organizer 2 and Wrye Bash support for modpack creators and modpack users. It will be possible to create modpack with Wrye Bash and install it with Mod Organizer 2.

- Steps support, it will be possible to define steps by modpack author which allow for:
  - running some additional tools in the middle of installation, running game to check that everything is ok, etc.,
  - have alternatives to current step (for example Performance vs Quality option),
  - ability to skip step.

- Semi-automated modpack creating - everyone should be able to make modpack without spending week doing so.

- Reasonable upgrade modpack procedure.

## Why?

Yes, there are other tools to do that, I was inspired by Automaton(TODO: link), but I find it not as flexible as I would like. And it is written in C#.

To answer why make such tool. Usually I spend so much time installing mods for my games, that I have no time to play. Tools as this will cut installation time to minimum. But what is VERY IMPORTANT for me - mod authors or Nexus are NOT HURT in the process. User still have to download mod himself. Only install process is automated. In fact, this not different from what mod managers are doing.

