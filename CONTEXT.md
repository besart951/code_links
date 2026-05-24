# Context

## Glossary

### Admin Activity Event

A time-ordered operational event shown to an Admin Actor. It can describe authentication activity, a security signal, a notification, an admin audit entry, or a runtime process log.

### Runtime Process Log

A file-backed Auth Service process log produced through Go `log.*`. Fatal process exits are persisted as `fatal:` entries so `/admin/logs` can show the crash after restart.
