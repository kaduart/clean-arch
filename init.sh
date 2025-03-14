rabbitmqctl delete_user guest
rabbitmqctl add_vhost /prod
rabbitmqctl set_permissions -p /prod admin ".*" ".*" ".*"
rabbitmqctl set_parameter -p /prod federation local_username 'admin'