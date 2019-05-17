Variables
{{ "hello" }}
{{ 'hello' }}
{{ "hell'o" }}

Filters
{{ 'Test'|slice:'1:3' }}
{{ '<div class=\"foo\"><ul class=\"foo\"><li class=\"foo\"><p class=\"foo\">This is a long test which will be cutted after some chars.</p></li></ul></div>'|truncatechars_html:25 }}
{{ '<a name="link"><p>This </a>is a long test which will be cutted after some chars.</p>'|truncatechars_html:25 }}

Tags
{% if 'Text' in complex.post %}text field in complex.post{% endif %}

Functions
{{ simple.func_variadic('hello') }}
