.PHONY: dispatch
dispatch:
	air -c .air.dispatch.toml

.PHONY: dispatch-debug
dispatch-debug:
	air -c .air.dispatch.debug.toml

.PHONY: notification
notification:
	air -c .air.notification.toml

.PHONY: notification-debug
notification-debug:
	air -c .air.notification.debug.toml

.PHONY: switchboard
switchboard:
	air -c .air.switchboard.toml

.PHONY: switchboard-debug
switchboard-debug:
	air -c .air.switchboard.debug.toml