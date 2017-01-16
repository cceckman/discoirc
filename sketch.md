If we skip the plugin system, what does the design look like, approximately?

* Common code
  * socket manager
    * Resolving directory, checking, cleaning up.
    * Health-checking?
  * logging
  * Client / service protos
    * "stream":
      * Network (e.g. notices)
      * Network #channel
      * Plugin (which *may* also be parameterized)
      * Needs to be, to some degree, stringable - for user to input.
      * Always allows input? Which may do nothing, e.g. listen-only plugins.
    * New-stream notification, for plugins.
* client:
  * Start UI.
  * Connect to server
    * Socket manager to list
    * Additional UI for socket manager: 
      * Selection
      * Cleaning up.
    * Start one if there isn't one running; then connect.
  * Get intended state:
    * Default
    * Last-good (last open)
    * command-line flags for what to open.
  * Request backfill & updates for stream.
    * Each stream is a different RPC stream? Probably easiest.
    * Scroll behavior? Don't want to pull the entire backfill.
      * So, view manager that separately pulls?
      * I don't know that that aligns with termui's behavior.
      * Maybe just ensure N lines above.
  * Able to start new windows! That's tricksy; need to know preferred termemu.
    * http://unix.stackexchange.com/questions/137782/launching-a-terminal-emulator-without-knowing-which-ones-are-installed:
      * i3-sensible-terminal if it's available
      * $TERMINAL (consistent with i3-sensible-terminal behavior)
      * xdg-terminal if it's available
      * x-terminal-emulator otherwise (Debian)
      * otherwise fallback to something.
* server:
  * Stream mananager
    * Connection manager for networks, channels
    * Plugin manager
    * Logging? Or is that effectively a plugin? Probably that.
  * Flow manager ?
    * Events propagating through streams ?
    * Once plugins are implemented, probably necessary.
  * View manager
    * Save, restore the last view.
  * Configuration manager
    * Load, reload, save configuration.
    * trigger, respond to updates in the stream manager.
    * trigger, respond to 
    * Saved config:
      * Default setup: autojoin, autoopen
      * Last-good setup: autojoin, autoopen
      * Default pref for what to open. 
