<?php
defined('BASEPATH') OR exit('No direct script access allowed');

/*
| -------------------------------------------------------------------
| DATABASE CONNECTIVITY SETTINGS
| -------------------------------------------------------------------
| This file will contain the settings needed to access your database.
|
| For complete instructions please consult the 'Database Connection'
| page of the User Guide.
|
| -------------------------------------------------------------------
| EXPLANATION OF VARIABLES
| -------------------------------------------------------------------

*/

$config = array();
$config['development']['socket_type'] = 'tcp'; //`tcp` or `unix`
$config['development']['socket'] = '/var/run/redis.sock'; // in case of `unix` socket type
$config['development']['host'] = '127.0.0.1';
$config['development']['password'] = 'test';
$config['development']['port'] = 6379;
$config['development']['timeout'] = 0;

$config['testing']['socket_type'] = 'tcp'; //`tcp` or `unix`
$config['testing']['socket'] = '/var/run/redis.sock'; // in case of `unix` socket type
$config['testing']['host'] = '127.0.0.1';
$config['testing']['password'] = 'test';
$config['testing']['port'] = 6379;
$config['testing']['timeout'] = 0;

$config['production']['socket_type'] = 'tcp'; //`tcp` or `unix`
$config['production']['socket'] = '/var/run/redis.sock'; // in case of `unix` socket type
$config['production']['host'] = '127.0.0.1';
$config['production']['password'] = 'test';
$config['production']['port'] = 6379;
$config['production']['timeout'] = 0;
