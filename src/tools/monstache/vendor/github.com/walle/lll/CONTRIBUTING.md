# Contributing

## Issues

When opening issues, try to be as clear as possible when describing the bug or
feature request. Tag the issue accordingly.

## Pull Requests

### Hack on lll

1. Install as usual (`go get github.com/walle/lll`)
2. Write code and tests for your new feature
3. Ensure everything works and the tests pass (see below)
4. Consider contributing your code upstream

### Contribute upstream

1. Fork ll on GitHub
2. Add your fork (`git remote add fork git@github.com:myuser/lll.git`)
3. Checkout your fork (`git checkout -t fork/master`)
4. Create your feature branch (`git checkout -b my-new-feature`)
5. Write code and tests for your new feature
6. Rebase against upstream to get changes \
(`git fetch origin && git rebase origin/master`)
7. Ensure everything works and the tests pass (see below)
8. Commit your changes
9. Push the branch to github (`git push fork my-new-feature`)
10. Create a new Pull Request on GitHub

Notice: Always use the original import path by installing with `go get`.

## Testing

To run the test suite use the command

```shell
$ go test -cover
```
