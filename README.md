git-hooks
========


Useful git-hooks in Golang

----------


ansible-lint
-------------

This hook checks [Ansible](https://www.ansible.com/) playbooks and/or roles for any syntax and best practices;  [ansible-lint](https://github.com/willthames/ansible-lint) for the same.

Additionally, the hook has some exceptions which can be tweaked as per requirement. 

#### Pre-Requisites

* go 1.8 or higher 

#### Build Executable

##### OS X

```sh
GOOS=darwin GOARCH=amd64 go build -o pre-commit pre-commit_ansible-lint.go
```

##### Windows

```sh
GOOS=windows GOARCH=386 go build -o pre-commit pre-commit_ansible-lint.go
```

**TIP**: More of these combinations can be looked using `go tool dist list`

#### Installation

Either place the built executable in `~/.git-templates/hooks/` directory so that the hook gets copied automagically each time a repo is cloned.

Or one can copy/move the executable directly in `~/.git/hooks/` directory. 

In Git GUI tools, only [SourceTree](https://www.sourcetreeapp.com/) has been tested on Windows and OSX.
