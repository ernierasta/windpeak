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
BTW: I was also looking into the Morrowind modding and while it has some quirks, it should be possible to support also Morrowind and especially OpenMW, where Linux support would be useful. I should collect some notes on this topic.

**2019-05-10**
Some investigation done and I had tested a lot of options how to detect changes. I need it to allow creators to create mod definitions (almost?) automatically if they are using Wrye Bash. For now there will be 3 processes. I will snapshot gamedir before mod is installed, then after installation/configuration. I will also make snapshot of unpacked mod. This will allow to make:

 * list of install actions (which files should be copied, where)
 * patch for ini files (mod configs, game configs)
 * binary patch!

Mod directory then will contain meta file with mod info and patches. Modpack author can also decide to put there files, which will overwrite destination file completely. And there are steps, but I have few ideas how to do them, we will see later.
This should be enough to create modpack. The only missing part I can think now is action of deleting files, which could be necessary for modpack updating (when you want remove some file or whole mod).

**2019-5-16**
No new code today. Some server work takes all my free time. But I started implementing 7z (zip & rar) unpacking ... and it was pain. There are currently 2 wrappers for 7zip functionality in golang. I wanted to go with library which uses 7zip.dll (7zip.so) which would be more elagant. While I figured out how to do that on Linux (you need to compile 2 c/c++ libraries) i had problems on Windows(not going into details, even if figure it out, there will be probably no option to generate 32-bit version). In short I will use 7z.exe (p7zip) wrapper probably. Need to make more tests. This is esential to have unpacking/packing done well.

BTW: I mentioned this project in youtube comments today. If you come from there - thanks for stopping by. :-)

## Planned features

- Multiplatform - Linux and Windows will work. Darwin will remain untested/unsupported until someone will want to test it there.

- Mod Organizer 2 and Wrye Bash support for modpack creators and modpack users. It will be possible to create modpack with Wrye Bash and install it with Mod Organizer 2.

- Steps support, it will be possible to define steps by modpack author which allow for:
  - running some additional tools in the middle of installation, running game to check that everything is ok, etc.,
  - have alternatives to current step (for example Performance vs Quality option),
  - ability to skip step.

- text files patches and binary patches. Yes, it will be possible to provide cleaned mod without distributing esm/esp file!

- Semi-automated modpack creating - everyone should be able to make modpack without spending week doing so.

- Reasonable upgrade modpack procedure.

## Why?

Yes, there are is [Automaton](https://github.com/metherul/Automaton) which is my inspiration, but I find it not as flexible as I would like. And it is written in C# so no option to me to contribute.

To answer - why I want to have such tool. Usually I spend so much time installing mods for my games, that I have no time to play. Tools as this will cut installation time to minimum. But what is VERY IMPORTANT for me - mod authors or Nexus are NOT HURT in the process. User still have to download mod himself. Only install process is automated. In fact, this not different from what mod managers are doing.

