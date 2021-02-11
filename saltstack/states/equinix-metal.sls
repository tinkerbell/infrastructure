{% for network in pillar['network']['addresses'] %}
{% if False == network['management'] %}
lo:{{ loop.index }}:
  network.managed:
    - enabled: True
    - type: eth
    - proto: static
    - ipaddr: {{ network['address'] }}
    - netmask: 255.255.255.255
{% endif %}
{% endfor %}
