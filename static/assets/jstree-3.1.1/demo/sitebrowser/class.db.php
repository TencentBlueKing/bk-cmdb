<?php
/*
// create instance protocol://user:pass@host/db?charset=UTF8&persist=TRUE&timezone=Europe/Sofia
$db = DB::get('mysqli://root@127.0.0.1/test');
$db = DB::get('oracle://root:pass@VIRT/?charset=AL32UTF8');
// execute a non-resulting query (returns boolean true)
$db->query('UPDATE table SET value = 1 WHERE id = ?', array(1));
// get all results as array
$db->all('SELECT * FROM table WHERE id = ?', array(1), "array_key", bool_skip_key, "assoc"/"num");
// get one result
$db->one('SELECT * FROM table WHERE id = ?', array(1), "assoc"/"num");
// get a traversable object to pass to foreach, or use count(), or use direct access: [INDEX]
$db->get('SELECT * FROM table WHERE id = ?', array(1), "assoc"/"num")[1];
*/

namespace
{
	class db
	{
		private function __construct() {
		}
		public function __clone() {
			throw new \vakata\database\Exception('Cannot clone static DB');
		}
		public static function get($settings = null) {
			return new \vakata\database\DBC($settings);
		}
		public static function getc($settings = null, \vakata\cache\ICache $c = null) {
			if($c === null) { $c = \vakata\cache\cache::inst(); }
			return new \vakata\database\DBCCached($settings, $c);
		}
	}
}

namespace vakata\database
{
	class Exception extends \Exception
	{
	}

	class Settings
	{
		public $type		= null;
		public $username	= 'root';
		public $password	= null;
		public $database	= null;
		public $servername	= 'localhost';
		public $serverport	= null;
		public $persist		= false;
		public $timezone	= null;
		public $charset		= 'UTF8';

		public function __construct($settings) {
			$str = parse_url($settings);
			if(!$str) {
				throw new Exception('Malformed DB settings string: ' . $settings);
			}
			if(array_key_exists('scheme',$str)) {
				$this->type			= rawurldecode($str['scheme']);
			}
			if(array_key_exists('user',$str)) {
				$this->username		= rawurldecode($str['user']);
			}
			if(array_key_exists('pass',$str)) {
				$this->password		= rawurldecode($str['pass']);
			}
			if(array_key_exists('path',$str)) {
				$this->database		= trim(rawurldecode($str['path']),'/');
			}
			if(array_key_exists('host',$str)) {
				$this->servername	= rawurldecode($str['host']);
			}
			if(array_key_exists('port',$str)) {
				$this->serverport	= rawurldecode($str['port']);
			}
			if(array_key_exists('query',$str)) {
				parse_str($str['query'], $str);
				$this->persist = (array_key_exists('persist', $str) && $str['persist'] === 'TRUE');
				if(array_key_exists('charset', $str)) {
					$this->charset = $str['charset'];
				}
				if(array_key_exists('timezone', $str)) {
					$this->timezone = $str['timezone'];
				}
			}
		}
	}

	interface IDB
	{
		public function connect();
		public function query($sql, $vars);
		public function get($sql, $data, $key, $skip_key, $mode);
		public function all($sql, $data, $key, $skip_key, $mode);
		public function one($sql, $data, $mode);
		public function raw($sql);
		public function prepare($sql);
		public function execute($data);
		public function disconnect();
	}

	interface IDriver
	{
		public function prepare($sql);
		public function execute($data);
		public function query($sql, $data);
		public function nextr($result);
		public function seek($result, $row);
		public function nf($result);
		public function af();
		public function insert_id();
		public function real_query($sql);
		public function get_settings();
	}

	abstract class ADriver implements IDriver
	{
		protected $lnk = null;
		protected $settings = null;

		public function __construct(Settings $settings) {
			$this->settings = $settings;
		}
		public function __destruct() {
			if($this->is_connected()) {
				$this->disconnect();
			}
		}
		public function get_settings() {
			return $this->settings;
		}

		public function connect() {
		}
		public function is_connected() {
			return $this->lnk !== null;
		}
		public function disconnect() {
		}
		public function query($sql, $data = array()) {
			return $this->execute($this->prepare($sql), $data);
		}
		public function prepare($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			return $sql;
		}
		public function execute($sql, $data = array()) {
			if(!$this->is_connected()) { $this->connect(); }
			if(!is_array($data)) { $data = array(); }
			$binder = '?';
			if(strpos($sql, $binder) !== false && is_array($data) && count($data)) {
				$tmp = explode($binder, $sql);
				if(!is_array($data)) { $data = array($data); }
				$data = array_values($data);
				if(count($data) >= count($tmp)) { $data = array_slice($data, 0, count($tmp)-1); }
				$sql = $tmp[0];
				foreach($data as $i => $v) {
					$sql .= $this->escape($v) . $tmp[($i + 1)];
				}
			}
			return $this->real_query($sql);
		}

		public function real_query($sql) {
			if(!$this->is_connected()) { $this->connect(); }
		}
		protected function escape($input) {
			if(is_array($input)) {
				foreach($input as $k => $v) {
					$input[$k] = $this->escape($v);
				}
				return implode(',',$input);
			}
			if(is_string($input)) {
				$input = addslashes($input);
				return "'".$input."'";
			}
			if(is_bool($input)) {
				return $input === false ? 0 : 1;
			}
			if(is_null($input)) {
				return 'NULL';
			}
			return $input;
		}

		public function nextr($result) {}
		public function nf($result) {}
		public function af() {}
		public function insert_id() {}
		public function seek($result, $row) {}
	}

	class Result implements \Iterator, \ArrayAccess, \Countable
	{
		protected $all  = null;
		protected $rdy  = false;
		protected $rslt	= null;
		protected $mode	= null;
		protected $fake	= null;
		protected $skip	= false;

		protected $fake_key	= 0;
		protected $real_key	= 0;
		public function __construct(Query $rslt, $key = null, $skip_key = false, $mode = 'assoc') {
			$this->rslt = $rslt;
			$this->mode = $mode;
			$this->fake = $key;
			$this->skip = $skip_key;
		}
		public function count() {
			return $this->rdy ? count($this->all) : $this->rslt->nf();
		}
		public function current() {
			if(!$this->count()) {
				return null;
			}
			if($this->rdy) {
				return current($this->all);
			}
			$tmp = $this->rslt->row();
			$row = array();
			switch($this->mode) {
				case 'num':
					foreach($tmp as $k => $v) {
						if(is_int($k)) {
							$row[$k] = $v;
						}
					}
					break;
				case 'both':
					$row = $tmp;
					break;
				case 'assoc':
				default:
					foreach($tmp as $k => $v) {
						if(!is_int($k)) {
							$row[$k] = $v;
						}
					}
					break;
			}
			if($this->fake) {
				$this->fake_key = $row[$this->fake];
			}
			if($this->skip) {
				unset($row[$this->fake]);
			}
			if(is_array($row) && count($row) === 1) {
				$row = current($row);
			}
			return $row;
		}
		public function key() {
			if($this->rdy) {
				return key($this->all);
			}
			return $this->fake ? $this->fake_key : $this->real_key;
		}
		public function next() {
			if($this->rdy) {
				return next($this->all);
			}
			$this->rslt->nextr();
			$this->real_key++;
		}
		public function rewind() {
			if($this->rdy) {
				return reset($this->all);
			}
			if($this->real_key !== null) {
				$this->rslt->seek(($this->real_key = 0));
			}
			$this->rslt->nextr();
		}
		public function valid() {
			if($this->rdy) {
				return current($this->all) !== false;
			}
			return $this->rslt->row() !== false && $this->rslt->row() !== null;
		}

		public function one() {
			$this->rewind();
			return $this->current();
		}
		public function get() {
			if(!$this->rdy) {
				$this->all = array();
				foreach($this as $k => $v) {
					$this->all[$k] = $v;
				}
				$this->rdy = true;
			}
			return $this->all;
		}
		public function offsetExists($offset) {
			if($this->rdy) {
				return isset($this->all[$offset]);
			}
			if($this->fake === null) {
				return $this->rslt->seek(($this->real_key = $offset));
			}
			$this->get();
			return isset($this->all[$offset]);
		}
		public function offsetGet($offset) {
			if($this->rdy) {
				return $this->all[$offset];
			}
			if($this->fake === null) {
				$this->rslt->seek(($this->real_key = $offset));
				$this->rslt->nextr();
				return $this->current();
			}
			$this->get();
			return $this->all[$offset];
		}
		public function offsetSet ($offset, $value ) {
			throw new Exception('Cannot set result');
		}
		public function offsetUnset ($offset) {
			throw new Exception('Cannot unset result');
		}
		public function __sleep() {
			$this->get();
			return array('all', 'rdy', 'mode', 'fake', 'skip');
		}
		public function __toString() {
			return print_r($this->get(), true);
		}
	}

	class Query
	{
		protected $drv = null;
		protected $sql = null;
		protected $prp = null;
		protected $rsl = null;
		protected $row = null;
		protected $num = null;
		protected $aff = null;
		protected $iid = null;

		public function __construct(IDriver $drv, $sql) {
			$this->drv = $drv;
			$this->sql = $sql;
			$this->prp = $this->drv->prepare($sql);
		}
		public function execute($vars = array()) {
			$this->rsl = $this->drv->execute($this->prp, $vars);
			$this->num = (is_object($this->rsl) || is_resource($this->rsl)) && is_callable(array($this->drv, 'nf')) ? (int)@$this->drv->nf($this->rsl) : 0;
			$this->aff = $this->drv->af();
			$this->iid = $this->drv->insert_id();
			return $this;
		}
		public function result($key = null, $skip_key = false, $mode = 'assoc') {
			return new Result($this, $key, $skip_key, $mode);
		}
		public function row() {
			return $this->row;
		}
		public function f($field) {
			return $this->row[$field];
		}
		public function nextr() {
			$this->row = $this->drv->nextr($this->rsl);
			return $this->row !== false && $this->row !== null;
		}
		public function seek($offset) {
			return @$this->drv->seek($this->rsl, $offset) ? true : false;
		}
		public function nf() {
			return $this->num;
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}
	}

	class DBC implements IDB
	{
		protected $drv = null;
		protected $que = null;

		public function __construct($drv = null) {
			if(!$drv && defined('DATABASE')) {
				$drv = DATABASE;
			}
			if(!$drv) {
				$this->error('Could not create database (no settings)');
			}
			if(is_string($drv)) {
				$drv = new \vakata\database\Settings($drv);
			}
			if($drv instanceof Settings) {
				$tmp = '\\vakata\\database\\' . $drv->type . '_driver';
				if(!class_exists($tmp)) {
					$this->error('Could not create database (no driver: '.$drv->type.')');
				}
				$drv = new $tmp($drv);
			}
			if(!($drv instanceof IDriver)) {
				$this->error('Could not create database - wrong driver');
			}
			$this->drv = $drv;
		}

		public function connect() {
			if(!$this->drv->is_connected()) {
				try {
					$this->drv->connect();
				}
				catch (Exception $e) {
					$this->error($e->getMessage(), 1);
				}
			}
			return true;
		}
		public function disconnect() {
			if($this->drv->is_connected()) {
				$this->drv->disconnect();
			}
		}

		public function prepare($sql) {
			try {
				$this->que = new Query($this->drv, $sql);
				return $this->que;
			} catch (Exception $e) {
				$this->error($e->getMessage(), 2);
			}
		}
		public function execute($data = array()) {
			try {
				return $this->que->execute($data);
			} catch (Exception $e) {
				$this->error($e->getMessage(), 3);
			}
		}
		public function query($sql, $data = array()) {
			try {
				$this->que = new Query($this->drv, $sql);
				return $this->que->execute($data);
			}
			catch (Exception $e) {
				$this->error($e->getMessage(), 4);
			}
		}
		public function get($sql, $data = array(), $key = null, $skip_key = false, $mode = 'assoc') {
			return $this->query($sql, $data)->result($key, $skip_key, $mode);
		}
		public function all($sql, $data = array(), $key = null, $skip_key = false, $mode = 'assoc') {
			return $this->get($sql, $data, $key, $skip_key, $mode)->get();
		}
		public function one($sql, $data = array(), $mode = 'assoc') {
			return $this->query($sql, $data)->result(null, false, $mode)->one();
		}
		public function raw($sql) {
			return $this->drv->real_query($sql);
		}
		public function get_driver() {
			return $this->drv->get_settings()->type;
		}

		public function __call($method, $args) {
			if($this->que && is_callable(array($this->que, $method))) {
				try {
					return call_user_func_array(array($this->que, $method), $args);
				} catch (Exception $e) {
					$this->error($e->getMessage(), 5);
				}
			}
		}

		protected final function error($error = '') {
			$dirnm = defined('LOGROOT') ? LOGROOT : realpath(dirname(__FILE__));
			@file_put_contents(
				$dirnm . DIRECTORY_SEPARATOR . '_errors_sql.log',
				'[' . date('d-M-Y H:i:s') . '] ' . $this->settings->type . ' > ' . preg_replace("@[\s\r\n\t]+@", ' ', $error) . "\n",
				FILE_APPEND
			);
			throw new Exception($error);
		}
	}

	class DBCCached extends DBC
	{
		protected $cache_inst = null;
		protected $cache_nmsp = null;
		public function __construct($settings = null, \vakata\cache\ICache $c = null) {
			parent::__construct($settings);
			$this->cache_inst = $c;
			$this->cache_nmsp = 'DBCCached_' . md5(serialize($this->drv->get_settings()));
		}
		public function cache($expires, $sql, $data = array(), $key = null, $skip_key = false, $mode = 'assoc') {
			$arg = func_get_args();
			array_shift($arg);
			$key = md5(serialize($arg));
			if(!$this->cache_inst) {
				return call_user_func_array(array($this, 'all'), $arg);
			}
			
			$tmp = $this->cache_inst->get($key, $this->cache_nmsp);
			if(!$tmp) {
				$this->cache_inst->prep($key, $this->cache_nmsp);
				$tmp = call_user_func_array(array($this, 'all'), $arg);
				$this->cache_inst->set($key, $tmp, $this->cache_nmsp, $expires);
			}
			return $tmp;
		}
		public function clear() {
			if($this->cache_inst) {
				$this->cache_inst->clear($this->cache_nmsp);
			}
		}
	}

	class mysqli_driver extends ADriver
	{
		protected $iid = 0;
		protected $aff = 0;
		protected $mnd = false;

		public function __construct($settings) {
			parent::__construct($settings);
			if(!$this->settings->serverport) { $this->settings->serverport = 3306; }
			$this->mnd = function_exists('mysqli_fetch_all');
		}

		public function connect() {
			$this->lnk = new \mysqli(
				($this->settings->persist ? 'p:' : '') . $this->settings->servername,
				$this->settings->username,
				$this->settings->password,
				$this->settings->database,
				$this->settings->serverport
			);
			if($this->lnk->connect_errno) {
				throw new Exception('Connect error: ' . $this->lnk->connect_errno);
			}
			if(!$this->lnk->set_charset($this->settings->charset)) {
				throw new Exception('Charset error: ' . $this->lnk->connect_errno);
			}
			if($this->settings->timezone) {
				@$this->lnk->query("SET time_zone = '" . addslashes($this->settings->timezone) . "'");
			}
			return true;
		}
		public function disconnect() {
			if($this->is_connected()) {
				@$this->lnk->close();
			}
		}
		public function real_query($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$temp = $this->lnk->query($sql);
			if(!$temp) {
				throw new Exception('Could not execute query : ' . $this->lnk->error . ' <'.$sql.'>');
			}
			$this->iid = $this->lnk->insert_id;
			$this->aff = $this->lnk->affected_rows;
			return $temp;
		}
		public function nextr($result) {
			if($this->mnd) {
				return $result->fetch_array(MYSQL_BOTH);
			}
			else {
				$ref = $result->result_metadata();
				if(!$ref) { return false; }
				$tmp = mysqli_fetch_fields($ref);
				if(!$tmp) { return false; }
				$ref = array();
				foreach($tmp as $col) { $ref[$col->name] = null; }
				$tmp = array();
				foreach($ref as $k => $v) { $tmp[] =& $ref[$k]; }
				if(!call_user_func_array(array($result, 'bind_result'), $tmp)) { return false; }
				if(!$result->fetch()) { return false; }
				$tmp = array();
				$i = 0;
				foreach($ref as $k => $v) { $tmp[$i++] = $v; $tmp[$k] = $v; }
				return $tmp;
			}
		}
		public function seek($result, $row) {
			$temp = $result->data_seek($row);
			return $temp;
		}
		public function nf($result) {
			return $result->num_rows;
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}
		public function prepare($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$temp = $this->lnk->prepare($sql);
			if(!$temp) {
				throw new Exception('Could not prepare : ' . $this->lnk->error . ' <'.$sql.'>');
			}
			return $temp;
		}
		public function execute($sql, $data = array()) {
			if(!$this->is_connected()) { $this->connect(); }
			if(!is_array($data)) { $data = array(); }
			if(is_string($sql)) {
				return parent::execute($sql, $data);
			}

			$data = array_values($data);
			if($sql->param_count) {
				if(count($data) < $sql->param_count) {
					throw new Exception('Prepared execute - not enough parameters.');
				}
				$ref = array('');
				foreach($data as $i => $v) {
					switch(gettype($v)) {
						case "boolean":
						case "integer":
							$data[$i] = (int)$v;
							$ref[0] .= 'i';
							$ref[$i+1] =& $data[$i];
							break;
						case "double":
							$ref[0] .= 'd';
							$ref[$i+1] =& $data[$i];
							break;
						case "array":
							$data[$i] = implode(',',$v);
							$ref[0] .= 's';
							$ref[$i+1] =& $data[$i];
							break;
						case "object":
						case "resource":
							$data[$i] = serialize($data[$i]);
							$ref[0] .= 's';
							$ref[$i+1] =& $data[$i];
							break;
						default:
							$ref[0] .= 's';
							$ref[$i+1] =& $data[$i];
							break;
					}
				}
				call_user_func_array(array($sql, 'bind_param'), $ref);
			}
			$rtrn = $sql->execute();
			if(!$this->mnd) {
				$sql->store_result();
			}
			if(!$rtrn) {
				throw new Exception('Prepared execute error : ' . $this->lnk->error);
			}
			$this->iid = $this->lnk->insert_id;
			$this->aff = $this->lnk->affected_rows;
			if(!$this->mnd) {
				return $sql->field_count ? $sql : $rtrn;
			}
			else {
				return $sql->field_count ? $sql->get_result() : $rtrn;
			}
		}

		protected function escape($input) {
			if(is_array($input)) {
				foreach($input as $k => $v) {
					$input[$k] = $this->escape($v);
				}
				return implode(',',$input);
			}
			if(is_string($input)) {
				$input = $this->lnk->real_escape_string($input);
				return "'".$input."'";
			}
			if(is_bool($input)) {
				return $input === false ? 0 : 1;
			}
			if(is_null($input)) {
				return 'NULL';
			}
			return $input;
		}
	}

	class mysql_driver extends ADriver
	{
		protected $iid = 0;
		protected $aff = 0;
		public function __construct($settings) {
			parent::__construct($settings);
			if(!$this->settings->serverport) { $this->settings->serverport = 3306; }
		}
		public function connect() {
			$this->lnk = ($this->settings->persist) ?
					@mysql_pconnect(
						$this->settings->servername.':'.$this->settings->serverport,
						$this->settings->username,
						$this->settings->password
					) :
					@mysql_connect(
						$this->settings->servername.':'.$this->settings->serverport,
						$this->settings->username,
						$this->settings->password
					);

			if($this->lnk === false || !mysql_select_db($this->settings->database, $this->lnk) || !mysql_query("SET NAMES '".$this->settings->charset."'", $this->lnk)) {
				throw new Exception('Connect error: ' . mysql_error());
			}
			if($this->settings->timezone) {
				@mysql_query("SET time_zone = '" . addslashes($this->settings->timezone) . "'", $this->lnk);
			}
			return true;
		}
		public function disconnect() {
			if(is_resource($this->lnk)) {
				mysql_close($this->lnk);
			}
		}

		public function real_query($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$temp = mysql_query($sql, $this->lnk);
			if(!$temp) {
				throw new Exception('Could not execute query : ' . mysql_error($this->lnk) . ' <'.$sql.'>');
			}
			$this->iid = mysql_insert_id($this->lnk);
			$this->aff = mysql_affected_rows($this->lnk);
			return $temp;
		}
		public function nextr($result) {
			return mysql_fetch_array($result, MYSQL_BOTH);
		}
		public function seek($result, $row) {
			$temp = @mysql_data_seek($result, $row);
			if(!$temp) {
				//throw new Exception('Could not seek : ' . mysql_error($this->lnk));
			}
			return $temp;
		}
		public function nf($result) {
			return mysql_num_rows($result);
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}

		protected function escape($input) {
			if(is_array($input)) {
				foreach($input as $k => $v) {
					$input[$k] = $this->escape($v);
				}
				return implode(',',$input);
			}
			if(is_string($input)) {
				$input = mysql_real_escape_string($input, $this->lnk);
				return "'".$input."'";
			}
			if(is_bool($input)) {
				return $input === false ? 0 : 1;
			}
			if(is_null($input)) {
				return 'NULL';
			}
			return $input;
		}
	}

	class postgre_driver extends ADriver
	{
		protected $iid = 0;
		protected $aff = 0;
		public function __construct($settings) {
			parent::__construct($settings);
			if(!$this->settings->serverport) { $this->settings->serverport = 5432; }
		}
		public function connect() {
			$this->lnk = ($this->settings->persist) ?
					@pg_pconnect(
						"host=" . $this->settings->servername . " " .
						"port=" . $this->settings->serverport . " " .
						"user=" . $this->settings->username . " " .
						"password=" . $this->settings->password . " " .
						"dbname=" . $this->settings->database . " " .
						"options='--client_encoding=".strtoupper($this->settings->charset)."' "
					) :
					@pg_connect(
						"host=" . $this->settings->servername . " " .
						"port=" . $this->settings->serverport . " " .
						"user=" . $this->settings->username . " " .
						"password=" . $this->settings->password . " " .
						"dbname=" . $this->settings->database . " " .
						"options='--client_encoding=".strtoupper($this->settings->charset)."' "
					);
			if($this->lnk === false) {
				throw new Exception('Connect error');
			}
			if($this->settings->timezone) {
				@pg_query($this->lnk, "SET TIME ZONE '".addslashes($this->settings->timezone)."'");
			}
			return true;
		}
		public function disconnect() {
			if(is_resource($this->lnk)) {
				pg_close($this->lnk);
			}
		}
		public function real_query($sql) {
			return $this->query($sql);
		}
		public function prepare($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$binder = '?';
			if(strpos($sql, $binder) !== false) {
				$tmp = explode($binder, $sql);
				$sql = $tmp[0];
				foreach($tmp as $i => $v) {
					$sql .= '$' . ($i + 1);
					if(isset($tmp[($i + 1)])) {
						$sql .= $tmp[($i + 1)];
					}
				}
			}
			return $sql;
		}
		public function execute($sql, $data = array()) {
			if(!$this->is_connected()) { $this->connect(); }
			if(!is_array($data)) { $data = array(); }
			$temp = (is_array($data) && count($data)) ? pg_query_params($this->lnk, $sql, $data) : pg_query_params($this->lnk, $sql, array());
			if(!$temp) {
				throw new Exception('Could not execute query : ' . pg_last_error($this->lnk) . ' <'.$sql.'>');
			}
			if(preg_match('@^\s*(INSERT|REPLACE)\s+INTO@i', $sql)) {
				$this->iid = pg_query($this->lnk, 'SELECT lastval()');
				$this->aff = pg_affected_rows($temp);
			}
			return $temp;
		}

		public function nextr($result) {
			return pg_fetch_array($result, NULL, PGSQL_BOTH);
		}
		public function seek($result, $row) {
			$temp = @pg_result_seek($result, $row);
			if(!$temp) {
				//throw new Exception('Could not seek : ' . pg_last_error($this->lnk));
			}
			return $temp;
		}
		public function nf($result) {
			return pg_num_rows($result);
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}

		// Функция mysql_query?
		//  - http://okbob.blogspot.com/2009/08/mysql-functions-for-postgresql.html
		//  - http://www.xach.com/aolserver/mysql-to-postgresql.html
		//  - REPLACE unixtimestamp / limit / curdate
	}

	class oracle_driver extends ADriver
	{
		protected $iid = 0;
		protected $aff = 0;

		public function connect() {
			$this->lnk = ($this->settings->persist) ?
					@oci_pconnect($this->settings->username, $this->settings->password, $this->settings->servername, $this->settings->charset) :
					@oci_connect ($this->settings->username, $this->settings->password, $this->settings->servername, $this->settings->charset);
			if($this->lnk === false) {
				throw new Exception('Connect error : ' . oci_error());
			}
			if($this->settings->timezone) {
				$this->real_query("ALTER session SET time_zone = '" . addslashes($this->settings->timezone) . "'");
			}
			return true;
		}
		public function disconnect() {
			if(is_resource($this->lnk)) {
				oci_close($this->lnk);
			}
		}
		public function real_query($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$temp = oci_parse($this->lnk, $sql);
			if(!$temp || !oci_execute($temp)) {
				throw new Exception('Could not execute real query : ' . oci_error($temp));
			}
			$this->aff = oci_num_rows($temp);
			return $temp;
		}

		public function prepare($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$binder = '?';
			if(strpos($sql, $binder) !== false && $vars !== false) {
				$tmp = explode($this->binder, $sql);
				$sql = $tmp[0];
				foreach($tmp as $i => $v) {
					$sql .= ':f' . $i;
					if(isset($tmp[($i + 1)])) {
						$sql .= $tmp[($i + 1)];
					}
				}
			}
			return oci_parse($this->lnk, $sql);
		}
		public function execute($sql, $data = array()) {
			if(!$this->is_connected()) { $this->connect(); }
			if(!is_array($data)) { $data = array(); }
			$data = array_values($data);
			foreach($data as $i => $v) {
				switch(gettype($v)) {
					case "boolean":
					case "integer":
						$data[$i] = (int)$v;
						oci_bind_by_name($sql, 'f'.$i, $data[$i], SQLT_INT);
						break;
					case "array":
						$data[$i] = implode(',',$v);
						oci_bind_by_name($sql, 'f'.$i, $data[$i]);
						break;
					case "object":
					case "resource":
						$data[$i] = serialize($data[$i]);
						oci_bind_by_name($sql, 'f'.$i, $data[$i]);
						break;
					default:
						oci_bind_by_name($sql, 'f'.$i, $data[$i]);
						break;
				}
			}
			$temp = oci_execute($sql);
			if(!$temp) {
				throw new Exception('Could not execute query : ' . oci_error($sql));
			}
			$this->aff = oci_num_rows($sql);

			/* TO DO: get iid
			if(!$seqname) { return $this->error('INSERT_ID not supported with no sequence.'); }
			$stm = oci_parse($this->link, 'SELECT '.strtoupper(str_replace("'",'',$seqname)).'.CURRVAL FROM DUAL');
			oci_execute($stm, $sql);
			$tmp = oci_fetch_array($stm);
			$tmp = $tmp[0];
			oci_free_statement($stm);
			*/
			return $sql;
		}
		public function nextr($result) {
			return oci_fetch_array($result, OCI_BOTH);
		}
		public function seek($result, $row) {
			$cnt = 0;
			while($cnt < $row) {
				if(oci_fetch_array($result, OCI_BOTH) === false) {
					return false;
				}
				$cnt++;
			}
			return true;
		}
		public function nf($result) {
			return oci_num_rows($result);
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}
	}

	class ibase_driver extends ADriver
	{
		protected $iid = 0;
		protected $aff = 0;
		public function __construct($settings) {
			parent::__construct($settings);
			if(!is_file($this->settings->database) && is_file('/'.$this->settings->database)) {
				$this->settings->database = '/'.$this->settings->database;
			}
			$this->settings->servername = ($this->settings->servername === 'localhost' || $this->settings->servername === '127.0.0.1' || $this->settings->servername === '') ?
				'' :
				$this->settings->servername . ':';
		}
		public function connect() {
			$this->lnk = ($this->settings->persist) ?
					@ibase_pconnect(
						$this->settings->servername . $this->settings->database,
						$this->settings->username,
						$this->settings->password,
						strtoupper($this->settings->charset)
					) :
					@ibase_connect(
						$this->settings->servername . $this->settings->database,
						$this->settings->username,
						$this->settings->password,
						strtoupper($this->settings->charset)
					);
			if($this->lnk === false) {
				throw new Exception('Connect error: ' . ibase_errmsg());
			}
			return true;
		}
		public function disconnect() {
			if(is_resource($this->lnk)) {
				ibase_close($this->lnk);
			}
		}

		public function real_query($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			$temp = ibase_query($sql, $this->lnk);
			if(!$temp) {
				throw new Exception('Could not execute query : ' . ibase_errmsg() . ' <'.$sql.'>');
			}
			//$this->iid = mysql_insert_id($this->lnk);
			$this->aff = ibase_affected_rows($this->lnk);
			return $temp;
		}
		public function prepare($sql) {
			if(!$this->is_connected()) { $this->connect(); }
			return ibase_prepare($this->lnk, $sql);
		}
		public function execute($sql, $data = array()) {
			if(!$this->is_connected()) { $this->connect(); }
			if(!is_array($data)) { $data = array(); }
			$data = array_values($data);
			foreach($data as $i => $v) {
				switch(gettype($v)) {
					case "boolean":
					case "integer":
						$data[$i] = (int)$v;
						break;
					case "array":
						$data[$i] = implode(',',$v);
						break;
					case "object":
					case "resource":
						$data[$i] = serialize($data[$i]);
						break;
				}
			}
			array_unshift($data, $sql);
			$temp = call_user_func_array("ibase_execute", $data);
			if(!$temp) {
				throw new Exception('Could not execute query : ' . ibase_errmsg() . ' <'.$sql.'>');
			}
			$this->aff = ibase_affected_rows($this->lnk);
			return $temp;
		}
		public function nextr($result) {
			return ibase_fetch_assoc($result, IBASE_TEXT);
		}
		public function seek($result, $row) {
			return false;
		}
		public function nf($result) {
			return false;
		}
		public function af() {
			return $this->aff;
		}
		public function insert_id() {
			return $this->iid;
		}
	}
}
