# Thoughts
These are my thoughts relating to this project.

* I hate school.
    * make the learning process smoother
* Make this an engaging learning experience
    * better peer <-> peer interaction
    * collaboration options
    * a nice showcase
* Needs to not be modelled after a teacher relationship
    * no "this is your task, this is what you do"
    * options
    * students can add their own tasks
        * this isn't for every class though...
        * plugin model?
* NEEDS TO BE FLUID
* Possible usage at hackEDU
    * easy applicable model
    * works with __people__
* I want myself to use it for managing classes
* I need it to not make me want to cry when i look at it
* Pluggable
* Motivated kids should be able to help other kids
* Encourage interactions on other forms of media
    * slack integration
    * facebook group intergration
    * twitter and instagram hash tagss
* Easy terminology for beginners to understand
* Wiki for sharing information
* Can we not use the word assignments...

## jarvis -- the code runner
jarvis is the code runner portion of steel. It is designed to be dynamic and be easily modifyable and playable by even new students.

* Handle ALL OF THE TEST SUITES
* There should be various optional ways of marking "assignments"
    * code length
    * code complexity
    * test passes
    * communication between hackers
* Easily extendable
* Pluggable just like everything else!
* It is important it is secure
    * hackers like trying to break things, it is an inherent issue :D
* Potentially have a network of `nodes` called ultron xD
* Needs to give nice output
    * Recognise errors which are being spewed out and how they are fixed
        * too complex for a multi-language system?
    * give tips to coders if they are coming across certain errors

## ultron -- the plugin system
ultron is the system which manages all of the integrations. An integration is a piece of code which adds *something* to the home bar. As well as this, they have access to the news feed and basic social functions. If they define the behaviour, interactions can be made inter-system.

* RPC model 
    * anything that supports RPC/JSON can be used to make a plugin
    * out of process?
* host your own
    * some sort of connectivity API
* use ours
    * we have built some plugins to manage how steel works. You can enable or disable them as you wish.
