# pkg

The pkg directory stores the contracts shared by every module: `transport` (binding, response envelopes, endpoint declarations), `apperr` (typed application errors), `event` (the bus interface), `registry` (the contribution conventions) and `job` (the scheduler contribution type). Packages here depend on nothing above them — the architecture test enforces it.
