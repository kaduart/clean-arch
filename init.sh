rabbitmqctl set_permissions -p /prod guest ".*" ".*" ".*"
rabbitmqctl set_parameter -p /prod federation local_username 'guest'