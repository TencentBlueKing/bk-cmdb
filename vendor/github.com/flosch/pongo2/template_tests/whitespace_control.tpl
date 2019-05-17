Standard whitespace control:
{% if true %}
Standard whitespace control
{% endif %}

Full Trim whitespace control:
{% if true -%}
Full Trim whitespace control
{%- endif %}

Useful with logic:
{%- if false %}
1st choice
{%- elif false %}
2nd choice
{%- elif true %}
3rd choice
{%- endif %}

Cycle without whitespace control:
{% for i in simple.multiple_item_list %}
{{ i }}
{% endfor %}

Cycle with whitespace control:
{% for i in simple.multiple_item_list %}
{{- i }}
{% endfor %}

Trim everything:
{% for i in simple.multiple_item_list -%}
{{ i }}
{%- endfor %}
