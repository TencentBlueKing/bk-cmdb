# Bootstrap Switch
[![Dependency Status](https://david-dm.org/nostalgiaz/bootstrap-switch.svg?style=flat)](https://david-dm.org/nostalgiaz/bootstrap-switch)
[![devDependency Status](https://david-dm.org/nostalgiaz/bootstrap-switch/dev-status.svg?style=flat)](https://david-dm.org/nostalgiaz/bootstrap-switch#info=devDependencies)
[![NPM Version](http://img.shields.io/npm/v/bootstrap-switch.svg?style=flat)](https://www.npmjs.org/)

Turn checkboxes and radio buttons into toggle switches.
This library is created by [Mattia Larentis](http://github.com/nostalgiaz) and maintained by the core team, with the help of the community.

To get started, check out [http://bootstrap-switch.org](http://bootstrap-switch.org)!

#### Core team
- [Mattia Larentis](http://github.com/nostalgiaz)
- [Emanuele Marchi](http://github.com/lostcrew)
- **you?** drop me a line.


## Demo and Documentation

- [Examples](http://www.bootstrap-switch.org/examples.html)
- [Options](http://www.bootstrap-switch.org/options.html)
- [Methods](http://www.bootstrap-switch.org/methods.html)
- [Events](http://www.bootstrap-switch.org/events.html)


## Getting started

Include the dependencies: jQuery, Bootstrap and Bootstrap Switch CSS + Javascript:

``` html
[...]
<link href="bootstrap.css" rel="stylesheet">
<link href="bootstrap-switch.css" rel="stylesheet">
<script src="jquery.js"></script>
<script src="bootstrap-switch.js"></script>
[...]
```

Add your checkbox:

```html
<input type="checkbox" name="my-checkbox" checked>
```

Initialize Bootstrap Switch on it:

```javascript
$("[name='my-checkbox']").bootstrapSwitch();
```

Enjoy.


## Supported browsers

IE9+ and all the other modern browsers.


## LESS + SASS

Import `src/less/bootstrap2/bootstrap-switch.less` for version <= 2.3.2 or `src/less/bootstrap3/bootstrap-switch.less` for version <= 3.3.4 in your compilation stack.


## Bugs and feature requests

Have a bug or a feature request? Please first search for existing and closed issues. If your problem or idea is not addressed yet, [please open a new issue](https://github.com/nostalgiaz/bootstrap-switch/issues/new). 

The new issue should contain both a summary of the issue and the browser/OS environment in which it occurs and a link to the playground you prefer with the reduced test case.
If suitable, include the steps required to reproduce the bug.

Please do not use the issue tracker for personal support requests: [Stack Overflow](http://stackoverflow.com/questions/tagged/bootstrap-switch) is a better place to get help.

#### Known issues

- Make sure `.form-control` is not applied to the input. Bootstrap does not support that, refer to [Checkboxes and radios](http://getbootstrap.com/css/#checkboxes-and-radios)


## Integrations

### AngularJs

Two custom directives are available:
- [angular-bootstrap-switch](https://github.com/frapontillo/angular-bootstrap-switch)
- [angular-toggle-switch](https://github.com/JumpLink/angular-toggle-switch)

### KnockoutJs

A Knockout binding handler is available [here](https://github.com/pauloortins/knockout-bootstrap-switch)

### NuGet

A NuGet package is available [here](https://github.com/blachniet/bootstrap-switch-nuget)


## License

Licensed under the MIT License
[https://github.com/nostalgiaz/bootstrap-switch/issues/347](https://github.com/nostalgiaz/bootstrap-switch/issues/347)

