/*global jQuery */
// wrap in IIFE and pass jQuery as $
(function ($, undefined) {
	"use strict";

	// some private plugin stuff if needed
	var private_var = null;

	// extending the defaults
	$.jstree.defaults.example = {
		example_option : 'example_val'
	};

	// the actual plugin code
	$.jstree.plugins.example = function (options, parent) {
		// own function
		this.example_function = function (arg) {
			// you can chain this method if needed and available
			if(parent.example_function) { parent.example_function.call(this, arg); }
		};

		// *SPECIAL* FUNCTIONS
		this.init = function (el, options) {
			// do not forget parent
			parent.init.call(this, el, options);
		};
		// bind events if needed
		this.bind = function () {
			// call parent function first
			parent.bind.call(this);
			// do(stuff);
		};
		// unbind events if needed (all in jquery namespace are taken care of by the core)
		this.unbind = function () {
			// do(stuff);
			// call parent function last
			parent.unbind.call(this);
		};
		this.teardown = function () {
			// do not forget parent
			parent.teardown.call(this);
		};
		// state management - get and restore
		this.get_state = function () {
			// always get state from parent first
			var state = parent.get_state.call(this);
			// add own stuff to state
			state.example = { 'var' : 'val' };
			return state;
		};
		this.set_state = function (state, callback) {
			// only process your part if parent returns true
			// there will be multiple times with false
			if(parent.set_state.call(this, state, callback)) {
				// check the key you set above
				if(state.example) {
					// do(stuff); // like calling this.example_function(state.example.var);
					// remove your part of the state, call again and RETURN FALSE, the next cycle will be TRUE
					delete state.example;
					this.set_state(state, callback);
					return false;
				}
				// return true if your state is gone (cleared in the previous step)
				return true;
			}
			// parent was false - return false too
			return false;
		};
		// node transportation
		this.get_json = function (obj, options, flat) {
			// get the node from the parent
			var tmp = parent.get_json.call(this, obj, options, flat), i, j;
			if($.isArray(tmp)) {
				for(i = 0, j = tmp.length; i < j; i++) {
					tmp[i].example = 'value';
				}
			}
			else {
				tmp.example = 'value';
			}
			// return the original / modified node
			return tmp;
		};
	};

	// attach to document ready if needed
	$(function () {
		// do(stuff);
	});

	// you can include the example plugin in all instances by default
	$.jstree.defaults.plugins.push("example");
})(jQuery);