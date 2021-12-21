# BUILDING

`make platforms VERSION=<semantic version>`

or

`make platforms BRANCH=<git branch> VERSION=<semantic version>`

or

`make platforms GO=<go version> BRANCH=<git branch> VERSION=<semantic version>`

Full-on release:

`make release` instead of `make platforms`

# TODO

# Notes

Under the hood, I rely on the Go Zenity package to display a save dialog:
- On MacOS, the actual dialog work is delegated to OSAScript.
- On Linux and Windows I may need to embed https://github.com/ncruces/zenity.

In fact, https://pkg.go.dev/github.com/gen2brain/dlgs looked like a fine package, until I fount out it can only be used to open files, not save them(!) ... plus, it also uses zenity under the hood.
