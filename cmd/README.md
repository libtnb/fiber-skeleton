# cmd

The cmd directory stores the entry point of each application, one directory per binary: `app` (the HTTP server), `cli` (management commands) and `gen` (the module generator). Entry points stay minimal — parse flags, assemble the injector, run.
