export default {
    'result': true,
    'bk_error_code': 0,
    'bk_error_msg': null,
    'data': {
        'unchanged': [
            {
                'process_template_id': 55,
                'process_template_name': 'mysqld',
                'service_instance_count': 77,
                'service_instances': [
                    {
                        'service_instance': {
                            'id': 122,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        }
                    }
                ]
            }
        ],
        'changed': [
            {
                'process_template_id': 56,
                'process_template_name': 'apache',
                'service_instance_count': 77,
                'service_instances': [
                    {
                        'service_instance': {
                            'id': 122,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        },
                        'changed_attributes': [
                            {
                                'property_id': 57,
                                'property_name': '端口',
                                'property_value': '2378',
                                'template_property_value': '2379'
                            }
                        ]
                    },
                    {
                        'service_instance': {
                            'id': 123,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        },
                        'changed_attributes': [
                            {
                                'property_id': 57,
                                'property_name': '端口',
                                'property_value': '2378',
                                'template_property_value': '2379'
                            }
                        ]
                    },
                    {
                        'service_instance': {
                            'id': 124,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        },
                        'changed_attributes': [
                            {
                                'property_id': 56,
                                'property_name': 'IP',
                                'property_value': '127.0.0.1',
                                'template_property_value': '0.0.0.0'
                            }
                        ]
                    }
                ]
            },
            {
                'process_template_id': 57,
                'process_template_name': 'apache',
                'service_instance_count': 77,
                'service_instances': [
                    {
                        'service_instance': {
                            'id': 125,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        },
                        'changed_attributes': [
                            {
                                'property_id': 57,
                                'property_name': '端口',
                                'property_value': '2378',
                                'template_property_value': '2379'
                            },
                            {
                                'property_id': 56,
                                'property_name': 'IP',
                                'property_value': '127.0.0.1',
                                'template_property_value': '0.0.0.0'
                            }
                        ]
                    }
                ]
            }
        ],
        'added': [
            {
                'process_template_id': 58,
                'process_template_name': 'tomcat',
                'service_instance_count': 77,
                'service_instances': [
                    {
                        'service_instance': {
                            'id': 126,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        }
                    }
                ]
            }
        ],
        'removed': [
            {
                'process_template_id': 59,
                'process_template_name': 'router',
                'service_instance_count': 77,
                'service_instances': [
                    {
                        'service_instance': {
                            'id': 127,
                            'name': '192.168.1.2_mysql_2379',
                            'service_template_id': 44,
                            'module_id': 66,
                            'host_id': 56
                        }
                    }
                ]
            }
        ]
    }
}
