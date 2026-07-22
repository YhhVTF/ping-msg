# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-07-21

### Changed

- Default window size is now bigger
- Focus will now automatically shift to the message entry
- Usernames on messages are now displayed smaller

### Fixed

- Bug where the client would freeze if sending a message was attempted while not connected to the server

## [0.1.0] - 2026-07-21

### Added

- Audio assets although they aren't being used yet
- Changelog
- Protocol for communicating with server
- Receiving singular messages from server
- Sending messages to server
- Username prompt

### Changed

- Text in the message entry is now sent to the server when submitted

### Fixed

- No longer connects to the server over and over
