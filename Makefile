.PHONY: dispatch
dispatch:
	air -build.args_bin=dispatch

.PHONY: dispatch-debug
dispatch-debug:
	air -c .air.debug.toml -build.args_bin=switchboard

.PHONY: notification
notification:
	air -build.args_bin=notification

.PHONY: notification-debug
notification-debug:
	air -c .air.debug.toml -build.args_bin=notification

.PHONY: switchboard
switchboard:
	air -build.args_bin=switchboard

.PHONY: switchboard-debug
switchboard-debug:
	air -c .air.debug.toml -build.args_bin=switchboard

.PHONY: web
web:
	air -build.args_bin=web

.PHONY: web-debug
web-debug:
	air -c .air.debug.toml -build.args_bin=web