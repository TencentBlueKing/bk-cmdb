<?php
// TO DO: better exceptions, use params
class tree
{
	protected $db = null;
	protected $options = null;
	protected $default = array(
		'structure_table'	=> 'structure',		// the structure table (containing the id, left, right, level, parent_id and position fields)
		'data_table'		=> 'structure',		// table for additional fields (apart from structure ones, can be the same as structure_table)
		'data2structure'	=> 'id',			// which field from the data table maps to the structure table
		'structure'			=> array(			// which field (value) maps to what in the structure (key)
			'id'			=> 'id',
			'left'			=> 'lft',
			'right'			=> 'rgt',
			'level'			=> 'lvl',
			'parent_id'		=> 'pid',
			'position'		=> 'pos'
		),
		'data'				=> array()			// array of additional fields from the data table
	);

	public function __construct(\vakata\database\IDB $db, array $options = array()) {
		$this->db = $db;
		$this->options = array_merge($this->default, $options);
	}

	public function get_node($id, $options = array()) {
		$node = $this->db->one("
			SELECT
				s.".implode(", s.", $this->options['structure']).",
				d.".implode(", d.", $this->options['data'])."
			FROM
				".$this->options['structure_table']." s,
				".$this->options['data_table']." d
			WHERE
				s.".$this->options['structure']['id']." = d.".$this->options['data2structure']." AND
				s.".$this->options['structure']['id']." = ".(int)$id
		);
		if(!$node) {
			throw new Exception('Node does not exist');
		}
		if(isset($options['with_children'])) {
			$node['children'] = $this->get_children($id, isset($options['deep_children']));
		}
		if(isset($options['with_path'])) {
			$node['path'] = $this->get_path($id);
		}
		return $node;
	}

	public function get_children($id, $recursive = false) {
		$sql = false;
		if($recursive) {
			$node = $this->get_node($id);
			$sql = "
				SELECT
					s.".implode(", s.", $this->options['structure']).",
					d.".implode(", d.", $this->options['data'])."
				FROM
					".$this->options['structure_table']." s,
					".$this->options['data_table']." d
				WHERE
					s.".$this->options['structure']['id']." = d.".$this->options['data2structure']." AND
					s.".$this->options['structure']['left']." > ".(int)$node[$this->options['structure']['left']]." AND
					s.".$this->options['structure']['right']." < ".(int)$node[$this->options['structure']['right']]."
				ORDER BY
					s.".$this->options['structure']['left']."
			";
		}
		else {
			$sql = "
				SELECT
					s.".implode(", s.", $this->options['structure']).",
					d.".implode(", d.", $this->options['data'])."
				FROM
					".$this->options['structure_table']." s,
					".$this->options['data_table']." d
				WHERE
					s.".$this->options['structure']['id']." = d.".$this->options['data2structure']." AND
					s.".$this->options['structure']['parent_id']." = ".(int)$id."
				ORDER BY
					s.".$this->options['structure']['position']."
			";
		}
		return $this->db->all($sql);
	}

	public function get_path($id) {
		$node = $this->get_node($id);
		$sql = false;
		if($node) {
			$sql = "
				SELECT
					s.".implode(", s.", $this->options['structure']).",
					d.".implode(", d.", $this->options['data'])."
				FROM
					".$this->options['structure_table']." s,
					".$this->options['data_table']." d
				WHERE
					s.".$this->options['structure']['id']." = d.".$this->options['data2structure']." AND
					s.".$this->options['structure']['left']." < ".(int)$node[$this->options['structure']['left']]." AND
					s.".$this->options['structure']['right']." > ".(int)$node[$this->options['structure']['right']]."
				ORDER BY
					s.".$this->options['structure']['left']."
			";
		}
		return $sql ? $this->db->all($sql) : false;
	}

	public function mk($parent, $position = 0, $data = array()) {
		$parent = (int)$parent;
		if($parent == 0) { throw new Exception('Parent is 0'); }
		$parent = $this->get_node($parent, array('with_children'=> true));
		if(!$parent['children']) { $position = 0; }
		if($parent['children'] && $position >= count($parent['children'])) { $position = count($parent['children']); }

		$sql = array();
		$par = array();

		// PREPARE NEW PARENT
		// update positions of all next elements
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$this->options['structure']["position"]." + 1
			WHERE
				".$this->options['structure']["parent_id"]." = ".(int)$parent[$this->options['structure']['id']]." AND
				".$this->options['structure']["position"]." >= ".$position."
			";
		$par[] = false;

		// update left indexes
		$ref_lft = false;
		if(!$parent['children']) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_lft = $parent['children'][(int)$position][$this->options['structure']["left"]];
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." + 2
			WHERE
				".$this->options['structure']["left"]." >= ".(int)$ref_lft."
			";
		$par[] = false;

		// update right indexes
		$ref_rgt = false;
		if(!$parent['children']) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_rgt = $parent['children'][(int)$position][$this->options['structure']["left"]] + 1;
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." + 2
			WHERE
				".$this->options['structure']["right"]." >= ".(int)$ref_rgt."
			";
		$par[] = false;

		// INSERT NEW NODE IN STRUCTURE
		$sql[] = "INSERT INTO ".$this->options['structure_table']." (".implode(",", $this->options['structure']).") VALUES (?".str_repeat(',?', count($this->options['structure']) - 1).")";
		$tmp = array();
		foreach($this->options['structure'] as $k => $v) {
			switch($k) {
				case 'id':
					$tmp[] = null;
					break;
				case 'left':
					$tmp[] = (int)$ref_lft;
					break;
				case 'right':
					$tmp[] = (int)$ref_lft + 1;
					break;
				case 'level':
					$tmp[] = (int)$parent[$v] + 1;
					break;
				case 'parent_id':
					$tmp[] = $parent[$this->options['structure']['id']];
					break;
				case 'position':
					$tmp[] = $position;
					break;
				default:
					$tmp[] = null;
			}
		}
		$par[] = $tmp;

		foreach($sql as $k => $v) {
			try {
				$this->db->query($v, $par[$k]);
			} catch(Exception $e) {
				$this->reconstruct();
				throw new Exception('Could not create');
			}
		}
		if($data && count($data)) {
			$node = $this->db->insert_id();
			if(!$this->rn($node,$data)) {
				$this->rm($node);
				throw new Exception('Could not rename after create');
			}
		}
		return $node;
	}

	public function mv($id, $parent, $position = 0) {
		$id			= (int)$id;
		$parent		= (int)$parent;
		if($parent == 0 || $id == 0 || $id == 1) {
			throw new Exception('Cannot move inside 0, or move root node');
		}

		$parent		= $this->get_node($parent, array('with_children'=> true, 'with_path' => true));
		$id			= $this->get_node($id, array('with_children'=> true, 'deep_children' => true, 'with_path' => true));
		if(!$parent['children']) {
			$position = 0;
		}
		if($id[$this->options['structure']['parent_id']] == $parent[$this->options['structure']['id']] && $position > $id[$this->options['structure']['position']]) {
			$position ++;
		}
		if($parent['children'] && $position >= count($parent['children'])) {
			$position = count($parent['children']);
		}
		if($id[$this->options['structure']['left']] < $parent[$this->options['structure']['left']] && $id[$this->options['structure']['right']] > $parent[$this->options['structure']['right']]) {
			throw new Exception('Could not move parent inside child');
		}

		$tmp = array();
		$tmp[] = (int)$id[$this->options['structure']["id"]];
		if($id['children'] && is_array($id['children'])) {
			foreach($id['children'] as $c) {
				$tmp[] = (int)$c[$this->options['structure']["id"]];
			}
		}
		$width = (int)$id[$this->options['structure']["right"]] - (int)$id[$this->options['structure']["left"]] + 1;

		$sql = array();

		// PREPARE NEW PARENT
		// update positions of all next elements
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$this->options['structure']["position"]." + 1
			WHERE
				".$this->options['structure']["id"]." != ".(int)$id[$this->options['structure']['id']]." AND
				".$this->options['structure']["parent_id"]." = ".(int)$parent[$this->options['structure']['id']]." AND
				".$this->options['structure']["position"]." >= ".$position."
			";

		// update left indexes
		$ref_lft = false;
		if(!$parent['children']) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_lft = $parent['children'][(int)$position][$this->options['structure']["left"]];
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." + ".$width."
			WHERE
				".$this->options['structure']["left"]." >= ".(int)$ref_lft." AND
				".$this->options['structure']["id"]." NOT IN(".implode(',',$tmp).")
			";
		// update right indexes
		$ref_rgt = false;
		if(!$parent['children']) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_rgt = $parent['children'][(int)$position][$this->options['structure']["left"]] + 1;
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." + ".$width."
			WHERE
				".$this->options['structure']["right"]." >= ".(int)$ref_rgt." AND
				".$this->options['structure']["id"]." NOT IN(".implode(',',$tmp).")
			";

		// MOVE THE ELEMENT AND CHILDREN
		// left, right and level
		$diff = $ref_lft - (int)$id[$this->options['structure']["left"]];

		if($diff > 0) { $diff = $diff - $width; }
		$ldiff = ((int)$parent[$this->options['structure']['level']] + 1) - (int)$id[$this->options['structure']['level']];
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." + ".$diff.",
					".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." + ".$diff.",
					".$this->options['structure']["level"]." = ".$this->options['structure']["level"]." + ".$ldiff."
				WHERE ".$this->options['structure']["id"]." IN(".implode(',',$tmp).")
		";
		// position and parent_id
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$position.",
					".$this->options['structure']["parent_id"]." = ".(int)$parent[$this->options['structure']["id"]]."
				WHERE ".$this->options['structure']["id"]."  = ".(int)$id[$this->options['structure']['id']]."
		";

		// CLEAN OLD PARENT
		// position of all next elements
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$this->options['structure']["position"]." - 1
			WHERE
				".$this->options['structure']["parent_id"]." = ".(int)$id[$this->options['structure']["parent_id"]]." AND
				".$this->options['structure']["position"]." > ".(int)$id[$this->options['structure']["position"]];
		// left indexes
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." - ".$width."
			WHERE
				".$this->options['structure']["left"]." > ".(int)$id[$this->options['structure']["right"]]." AND
				".$this->options['structure']["id"]." NOT IN(".implode(',',$tmp).")
		";
		// right indexes
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." - ".$width."
			WHERE
				".$this->options['structure']["right"]." > ".(int)$id[$this->options['structure']["right"]]." AND
				".$this->options['structure']["id"]." NOT IN(".implode(',',$tmp).")
		";

		foreach($sql as $k => $v) {
			//echo preg_replace('@[\s\t]+@',' ',$v) ."\n";
			try {
				$this->db->query($v);
			} catch(Exception $e) {
				$this->reconstruct();
				throw new Exception('Error moving');
			}
		}
		return true;
	}

	public function cp($id, $parent, $position = 0) {
		$id			= (int)$id;
		$parent		= (int)$parent;
		if($parent == 0 || $id == 0 || $id == 1) {
			throw new Exception('Could not copy inside parent 0, or copy root nodes');
		}

		$parent		= $this->get_node($parent, array('with_children'=> true, 'with_path' => true));
		$id			= $this->get_node($id, array('with_children'=> true, 'deep_children' => true, 'with_path' => true));
		$old_nodes	= $this->db->get("
			SELECT * FROM ".$this->options['structure_table']."
			WHERE ".$this->options['structure']["left"]." > ".$id[$this->options['structure']["left"]]." AND ".$this->options['structure']["right"]." < ".$id[$this->options['structure']["right"]]."
			ORDER BY ".$this->options['structure']["left"]."
		");
		if(!$parent['children']) {
			$position = 0;
		}
		if($id[$this->options['structure']['parent_id']] == $parent[$this->options['structure']['id']] && $position > $id[$this->options['structure']['position']]) {
			//$position ++;
		}
		if($parent['children'] && $position >= count($parent['children'])) {
			$position = count($parent['children']);
		}

		$tmp = array();
		$tmp[] = (int)$id[$this->options['structure']["id"]];
		if($id['children'] && is_array($id['children'])) {
			foreach($id['children'] as $c) {
				$tmp[] = (int)$c[$this->options['structure']["id"]];
			}
		}
		$width = (int)$id[$this->options['structure']["right"]] - (int)$id[$this->options['structure']["left"]] + 1;

		$sql = array();

		// PREPARE NEW PARENT
		// update positions of all next elements
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$this->options['structure']["position"]." + 1
			WHERE
				".$this->options['structure']["parent_id"]." = ".(int)$parent[$this->options['structure']['id']]." AND
				".$this->options['structure']["position"]." >= ".$position."
			";

		// update left indexes
		$ref_lft = false;
		if(!$parent['children']) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_lft = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_lft = $parent['children'][(int)$position][$this->options['structure']["left"]];
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." + ".$width."
			WHERE
				".$this->options['structure']["left"]." >= ".(int)$ref_lft."
			";
		// update right indexes
		$ref_rgt = false;
		if(!$parent['children']) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else if(!isset($parent['children'][$position])) {
			$ref_rgt = $parent[$this->options['structure']["right"]];
		}
		else {
			$ref_rgt = $parent['children'][(int)$position][$this->options['structure']["left"]] + 1;
		}
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." + ".$width."
			WHERE
				".$this->options['structure']["right"]." >= ".(int)$ref_rgt."
			";

		// MOVE THE ELEMENT AND CHILDREN
		// left, right and level
		$diff = $ref_lft - (int)$id[$this->options['structure']["left"]];

		if($diff <= 0) { $diff = $diff - $width; }
		$ldiff = ((int)$parent[$this->options['structure']['level']] + 1) - (int)$id[$this->options['structure']['level']];

		// build all fields + data table
		$fields = array_combine($this->options['structure'], $this->options['structure']);
		unset($fields['id']);
		$fields[$this->options['structure']["left"]] = $this->options['structure']["left"]." + ".$diff;
		$fields[$this->options['structure']["right"]] = $this->options['structure']["right"]." + ".$diff;
		$fields[$this->options['structure']["level"]] = $this->options['structure']["level"]." + ".$ldiff;
		$sql[] = "
			INSERT INTO ".$this->options['structure_table']." ( ".implode(',',array_keys($fields))." )
			SELECT ".implode(',',array_values($fields))." FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']["id"]." IN (".implode(",", $tmp).")
			ORDER BY ".$this->options['structure']["level"]." ASC";

		foreach($sql as $k => $v) {
			try {
				$this->db->query($v);
			} catch(Exception $e) {
				$this->reconstruct();
				throw new Exception('Error copying');
			}
		}
		$iid = (int)$this->db->insert_id();

		try {
			$this->db->query("
				UPDATE ".$this->options['structure_table']."
					SET ".$this->options['structure']["position"]." = ".$position.",
						".$this->options['structure']["parent_id"]." = ".(int)$parent[$this->options['structure']["id"]]."
					WHERE ".$this->options['structure']["id"]."  = ".$iid."
			");
		} catch(Exception $e) {
			$this->rm($iid);
			$this->reconstruct();
			throw new Exception('Could not update adjacency after copy');
		}
		$fields = $this->options['data'];
		unset($fields['id']);
		$update_fields = array();
		foreach($fields as $f) {
			$update_fields[] = $f.'=VALUES('.$f.')';
		}
		$update_fields = implode(',', $update_fields);
		if(count($fields)) {
			try {
				$this->db->query("
						INSERT INTO ".$this->options['data_table']." (".$this->options['data2structure'].",".implode(",",$fields).")
						SELECT ".$iid.",".implode(",",$fields)." FROM ".$this->options['data_table']." WHERE ".$this->options['data2structure']." = ".$id[$this->options['data2structure']]."
						ON DUPLICATE KEY UPDATE ".$update_fields."
				");
			}
			catch(Exception $e) {
				$this->rm($iid);
				$this->reconstruct();
				throw new Exception('Could not update data after copy');
			}
		}

		// manually fix all parent_ids and copy all data
		$new_nodes = $this->db->get("
			SELECT * FROM ".$this->options['structure_table']."
			WHERE ".$this->options['structure']["left"]." > ".$ref_lft." AND ".$this->options['structure']["right"]." < ".($ref_lft + $width - 1)." AND ".$this->options['structure']["id"]." != ".$iid."
			ORDER BY ".$this->options['structure']["left"]."
		");
		$parents = array();
		foreach($new_nodes as $node) {
			if(!isset($parents[$node[$this->options['structure']["left"]]])) { $parents[$node[$this->options['structure']["left"]]] = $iid; }
			for($i = $node[$this->options['structure']["left"]] + 1; $i < $node[$this->options['structure']["right"]]; $i++) {
				$parents[$i] = $node[$this->options['structure']["id"]];
			}
		}
		$sql = array();
		foreach($new_nodes as $k => $node) {
			$sql[] = "
				UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["parent_id"]." = ".$parents[$node[$this->options['structure']["left"]]]."
				WHERE ".$this->options['structure']["id"]." = ".(int)$node[$this->options['structure']["id"]]."
			";
			if(count($fields)) {
				$up = "";
				foreach($fields as $f)
				$sql[] = "
					INSERT INTO ".$this->options['data_table']." (".$this->options['data2structure'].",".implode(",",$fields).")
					SELECT ".(int)$node[$this->options['structure']["id"]].",".implode(",",$fields)." FROM ".$this->options['data_table']."
						WHERE ".$this->options['data2structure']." = ".$old_nodes[$k][$this->options['structure']['id']]."
					ON DUPLICATE KEY UPDATE ".$update_fields."
				";
			}
		}
		//var_dump($sql);
		foreach($sql as $k => $v) {
			try {
				$this->db->query($v);
			} catch(Exception $e) {
				$this->rm($iid);
				$this->reconstruct();
				throw new Exception('Error copying');
			}
		}
		return $iid;
	}

	public function rm($id) {
		$id = (int)$id;
		if(!$id || $id === 1) { throw new Exception('Could not create inside roots'); }
		$data = $this->get_node($id, array('with_children' => true, 'deep_children' => true));
		$lft = (int)$data[$this->options['structure']["left"]];
		$rgt = (int)$data[$this->options['structure']["right"]];
		$pid = (int)$data[$this->options['structure']["parent_id"]];
		$pos = (int)$data[$this->options['structure']["position"]];
		$dif = $rgt - $lft + 1;

		$sql = array();
		// deleting node and its children from structure
		$sql[] = "
			DELETE FROM ".$this->options['structure_table']."
			WHERE ".$this->options['structure']["left"]." >= ".(int)$lft." AND ".$this->options['structure']["right"]." <= ".(int)$rgt."
		";
		// shift left indexes of nodes right of the node
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["left"]." = ".$this->options['structure']["left"]." - ".(int)$dif."
			WHERE ".$this->options['structure']["left"]." > ".(int)$rgt."
		";
		// shift right indexes of nodes right of the node and the node's parents
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["right"]." = ".$this->options['structure']["right"]." - ".(int)$dif."
			WHERE ".$this->options['structure']["right"]." > ".(int)$lft."
		";
		// Update position of siblings below the deleted node
		$sql[] = "
			UPDATE ".$this->options['structure_table']."
				SET ".$this->options['structure']["position"]." = ".$this->options['structure']["position"]." - 1
			WHERE ".$this->options['structure']["parent_id"]." = ".$pid." AND ".$this->options['structure']["position"]." > ".(int)$pos."
		";
		// delete from data table
		if($this->options['data_table']) {
			$tmp = array();
			$tmp[] = (int)$data['id'];
			if($data['children'] && is_array($data['children'])) {
				foreach($data['children'] as $v) {
					$tmp[] = (int)$v['id'];
				}
			}
			$sql[] = "DELETE FROM ".$this->options['data_table']." WHERE ".$this->options['data2structure']." IN (".implode(',',$tmp).")";
		}

		foreach($sql as $v) {
			try {
				$this->db->query($v);
			} catch(Exception $e) {
				$this->reconstruct();
				throw new Exception('Could not remove');
			}
		}
		return true;
	}

	public function rn($id, $data) {
		if(!(int)$this->db->one('SELECT 1 AS res FROM '.$this->options['structure_table'].' WHERE '.$this->options['structure']['id'].' = '.(int)$id)) {
			throw new Exception('Could not rename non-existing node');
		}
		$tmp = array();
		foreach($this->options['data'] as $v) {
			if(isset($data[$v])) {
				$tmp[$v] = $data[$v];
			}
		}
		if(count($tmp)) {
			$tmp[$this->options['data2structure']] = $id;
			$sql = "
				INSERT INTO
					".$this->options['data_table']." (".implode(',', array_keys($tmp)).")
					VALUES(?".str_repeat(',?', count($tmp) - 1).")
				ON DUPLICATE KEY UPDATE
					".implode(' = ?, ', array_keys($tmp))." = ?";
			$par = array_merge(array_values($tmp), array_values($tmp));
			try {
				$this->db->query($sql, $par);
			}
			catch(Exception $e) {
				throw new Exception('Could not rename');
			}
		}
		return true;
	}

	public function analyze($get_errors = false) {
		$report = array();
		if((int)$this->db->one("SELECT COUNT(".$this->options['structure']["id"].") AS res FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']["parent_id"]." = 0") !== 1) {
			$report[] = "No or more than one root node.";
		}
		if((int)$this->db->one("SELECT ".$this->options['structure']["left"]." AS res FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']["parent_id"]." = 0") !== 1) {
			$report[] = "Root node's left index is not 1.";
		}
		if((int)$this->db->one("
			SELECT
				COUNT(".$this->options['structure']['id'].") AS res
			FROM ".$this->options['structure_table']." s
			WHERE
				".$this->options['structure']["parent_id"]." != 0 AND
				(SELECT COUNT(".$this->options['structure']['id'].") FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']["id"]." = s.".$this->options['structure']["parent_id"].") = 0") > 0
		) {
			$report[] = "Missing parents.";
		}
		if(
			(int)$this->db->one("SELECT MAX(".$this->options['structure']["right"].") AS res FROM ".$this->options['structure_table']) / 2 !=
			(int)$this->db->one("SELECT COUNT(".$this->options['structure']["id"].") AS res FROM ".$this->options['structure_table'])
		) {
			$report[] = "Right index does not match node count.";
		}
		if(
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["right"].") AS res FROM ".$this->options['structure_table']) !=
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["left"].") AS res FROM ".$this->options['structure_table'])
		) {
			$report[] = "Duplicates in nested set.";
		}
		if(
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["id"].") AS res FROM ".$this->options['structure_table']) !=
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["left"].") AS res FROM ".$this->options['structure_table'])
		) {
			$report[] = "Left indexes not unique.";
		}
		if(
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["id"].") AS res FROM ".$this->options['structure_table']) !=
			(int)$this->db->one("SELECT COUNT(DISTINCT ".$this->options['structure']["right"].") AS res FROM ".$this->options['structure_table'])
		) {
			$report[] = "Right indexes not unique.";
		}
		if(
			(int)$this->db->one("
				SELECT
					s1.".$this->options['structure']["id"]." AS res
				FROM ".$this->options['structure_table']." s1, ".$this->options['structure_table']." s2
				WHERE
					s1.".$this->options['structure']['id']." != s2.".$this->options['structure']['id']." AND
					s1.".$this->options['structure']['left']." = s2.".$this->options['structure']['right']."
				LIMIT 1")
		) {
			$report[] = "Nested set - matching left and right indexes.";
		}
		if(
			(int)$this->db->one("
				SELECT
					".$this->options['structure']["id"]." AS res
				FROM ".$this->options['structure_table']." s
				WHERE
					".$this->options['structure']['position']." >= (
						SELECT
							COUNT(".$this->options['structure']["id"].")
						FROM ".$this->options['structure_table']."
						WHERE ".$this->options['structure']['parent_id']." = s.".$this->options['structure']['parent_id']."
					)
				LIMIT 1") ||
			(int)$this->db->one("
				SELECT
					s1.".$this->options['structure']["id"]." AS res
				FROM ".$this->options['structure_table']." s1, ".$this->options['structure_table']." s2
				WHERE
					s1.".$this->options['structure']['id']." != s2.".$this->options['structure']['id']." AND
					s1.".$this->options['structure']['parent_id']." = s2.".$this->options['structure']['parent_id']." AND
					s1.".$this->options['structure']['position']." = s2.".$this->options['structure']['position']."
				LIMIT 1")
		) {
			$report[] = "Positions not correct.";
		}
		if((int)$this->db->one("
			SELECT
				COUNT(".$this->options['structure']["id"].") FROM ".$this->options['structure_table']." s
			WHERE
				(
					SELECT
						COUNT(".$this->options['structure']["id"].")
					FROM ".$this->options['structure_table']."
					WHERE
						".$this->options['structure']["right"]." < s.".$this->options['structure']["right"]." AND
						".$this->options['structure']["left"]." > s.".$this->options['structure']["left"]." AND
						".$this->options['structure']["level"]." = s.".$this->options['structure']["level"]." + 1
				) !=
				(
					SELECT
						COUNT(*)
					FROM ".$this->options['structure_table']."
					WHERE
						".$this->options['structure']["parent_id"]." = s.".$this->options['structure']["id"]."
				)")
		) {
			$report[] = "Adjacency and nested set do not match.";
		}
		if(
			$this->options['data_table'] &&
			(int)$this->db->one("
				SELECT
					COUNT(".$this->options['structure']["id"].") AS res
				FROM ".$this->options['structure_table']." s
				WHERE
					(SELECT COUNT(".$this->options['data2structure'].") FROM ".$this->options['data_table']." WHERE ".$this->options['data2structure']." = s.".$this->options['structure']["id"].") = 0
			")
		) {
			$report[] = "Missing records in data table.";
		}
		if(
			$this->options['data_table'] &&
			(int)$this->db->one("
				SELECT
					COUNT(".$this->options['data2structure'].") AS res
				FROM ".$this->options['data_table']." s
				WHERE
					(SELECT COUNT(".$this->options['structure']["id"].") FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']["id"]." = s.".$this->options['data2structure'].") = 0
			")
		) {
			$report[] = "Dangling records in data table.";
		}
		return $get_errors ? $report : count($report) == 0;
	}

	public function reconstruct($analyze = true) {
		if($analyze && $this->analyze()) { return true; }

		if(!$this->db->query("" .
			"CREATE TEMPORARY TABLE temp_tree (" .
				"".$this->options['structure']["id"]." INTEGER NOT NULL, " .
				"".$this->options['structure']["parent_id"]." INTEGER NOT NULL, " .
				"". $this->options['structure']["position"]." INTEGER NOT NULL" .
			") "
		)) { return false; }
		if(!$this->db->query("" .
			"INSERT INTO temp_tree " .
				"SELECT " .
					"".$this->options['structure']["id"].", " .
					"".$this->options['structure']["parent_id"].", " .
					"".$this->options['structure']["position"]." " .
				"FROM ".$this->options['structure_table'].""
		)) { return false; }

		if(!$this->db->query("" .
			"CREATE TEMPORARY TABLE temp_stack (" .
				"".$this->options['structure']["id"]." INTEGER NOT NULL, " .
				"".$this->options['structure']["left"]." INTEGER, " .
				"".$this->options['structure']["right"]." INTEGER, " .
				"".$this->options['structure']["level"]." INTEGER, " .
				"stack_top INTEGER NOT NULL, " .
				"".$this->options['structure']["parent_id"]." INTEGER, " .
				"".$this->options['structure']["position"]." INTEGER " .
			") "
		)) { return false; }

		$counter = 2;
		if(!$this->db->query("SELECT COUNT(*) FROM temp_tree")) {
			return false;
		}
		$this->db->nextr();
		$maxcounter = (int) $this->db->f(0) * 2;
		$currenttop = 1;
		if(!$this->db->query("" .
			"INSERT INTO temp_stack " .
				"SELECT " .
					"".$this->options['structure']["id"].", " .
					"1, " .
					"NULL, " .
					"0, " .
					"1, " .
					"".$this->options['structure']["parent_id"].", " .
					"".$this->options['structure']["position"]." " .
				"FROM temp_tree " .
				"WHERE ".$this->options['structure']["parent_id"]." = 0"
		)) { return false; }
		if(!$this->db->query("DELETE FROM temp_tree WHERE ".$this->options['structure']["parent_id"]." = 0")) {
			return false;
		}

		while ($counter <= $maxcounter) {
			if(!$this->db->query("" .
				"SELECT " .
					"temp_tree.".$this->options['structure']["id"]." AS tempmin, " .
					"temp_tree.".$this->options['structure']["parent_id"]." AS pid, " .
					"temp_tree.".$this->options['structure']["position"]." AS lid " .
				"FROM temp_stack, temp_tree " .
				"WHERE " .
					"temp_stack.".$this->options['structure']["id"]." = temp_tree.".$this->options['structure']["parent_id"]." AND " .
					"temp_stack.stack_top = ".$currenttop." " .
				"ORDER BY temp_tree.".$this->options['structure']["position"]." ASC LIMIT 1"
			)) { return false; }

			if($this->db->nextr()) {
				$tmp = $this->db->f("tempmin");

				$q = "INSERT INTO temp_stack (stack_top, ".$this->options['structure']["id"].", ".$this->options['structure']["left"].", ".$this->options['structure']["right"].", ".$this->options['structure']["level"].", ".$this->options['structure']["parent_id"].", ".$this->options['structure']["position"].") VALUES(".($currenttop + 1).", ".$tmp.", ".$counter.", NULL, ".$currenttop.", ".$this->db->f("pid").", ".$this->db->f("lid").")";
				if(!$this->db->query($q)) {
					return false;
				}
				if(!$this->db->query("DELETE FROM temp_tree WHERE ".$this->options['structure']["id"]." = ".$tmp)) {
					return false;
				}
				$counter++;
				$currenttop++;
			}
			else {
				if(!$this->db->query("" .
					"UPDATE temp_stack SET " .
						"".$this->options['structure']["right"]." = ".$counter.", " .
						"stack_top = -stack_top " .
					"WHERE stack_top = ".$currenttop
				)) { return false; }
				$counter++;
				$currenttop--;
			}
		}

		$temp_fields = $this->options['structure'];
		unset($temp_fields["parent_id"]);
		unset($temp_fields["position"]);
		unset($temp_fields["left"]);
		unset($temp_fields["right"]);
		unset($temp_fields["level"]);
		if(count($temp_fields) > 1) {
			if(!$this->db->query("" .
				"CREATE TEMPORARY TABLE temp_tree2 " .
					"SELECT ".implode(", ", $temp_fields)." FROM ".$this->options['structure_table']." "
			)) { return false; }
		}
		if(!$this->db->query("TRUNCATE TABLE ".$this->options['structure_table']."")) {
			return false;
		}
		if(!$this->db->query("" .
			"INSERT INTO ".$this->options['structure_table']." (" .
					"".$this->options['structure']["id"].", " .
					"".$this->options['structure']["parent_id"].", " .
					"".$this->options['structure']["position"].", " .
					"".$this->options['structure']["left"].", " .
					"".$this->options['structure']["right"].", " .
					"".$this->options['structure']["level"]." " .
				") " .
				"SELECT " .
					"".$this->options['structure']["id"].", " .
					"".$this->options['structure']["parent_id"].", " .
					"".$this->options['structure']["position"].", " .
					"".$this->options['structure']["left"].", " .
					"".$this->options['structure']["right"].", " .
					"".$this->options['structure']["level"]." " .
				"FROM temp_stack " .
				"ORDER BY ".$this->options['structure']["id"].""
		)) {
			return false;
		}
		if(count($temp_fields) > 1) {
			$sql = "" .
				"UPDATE ".$this->options['structure_table']." v, temp_tree2 SET v.".$this->options['structure']["id"]." = v.".$this->options['structure']["id"]." ";
			foreach($temp_fields as $k => $v) {
				if($k == "id") continue;
				$sql .= ", v.".$v." = temp_tree2.".$v." ";
			}
			$sql .= " WHERE v.".$this->options['structure']["id"]." = temp_tree2.".$this->options['structure']["id"]." ";
			if(!$this->db->query($sql)) {
				return false;
			}
		}
		// fix positions
		$nodes = $this->db->get("SELECT ".$this->options['structure']['id'].", ".$this->options['structure']['parent_id']." FROM ".$this->options['structure_table']." ORDER BY ".$this->options['structure']['parent_id'].", ".$this->options['structure']['position']);
		$last_parent = false;
		$last_position = false;
		foreach($nodes as $node) {
			if((int)$node[$this->options['structure']['parent_id']] !== $last_parent) {
				$last_position = 0;
				$last_parent = (int)$node[$this->options['structure']['parent_id']];
			}
			$this->db->query("UPDATE ".$this->options['structure_table']." SET ".$this->options['structure']['position']." = ".$last_position." WHERE ".$this->options['structure']['id']." = ".(int)$node[$this->options['structure']['id']]);
			$last_position++;
		}
		if($this->options['data_table'] != $this->options['structure_table']) {
			// fix missing data records
			$this->db->query("
				INSERT INTO
					".$this->options['data_table']." (".implode(',',$this->options['data']).")
				SELECT ".$this->options['structure']['id']." ".str_repeat(", ".$this->options['structure']['id'], count($this->options['data']) - 1)."
				FROM ".$this->options['structure_table']." s
				WHERE (SELECT COUNT(".$this->options['data2structure'].") FROM ".$this->options['data_table']." WHERE ".$this->options['data2structure']." = s.".$this->options['structure']['id'].") = 0 "
			);
			// remove dangling data records
			$this->db->query("
				DELETE FROM
					".$this->options['data_table']."
				WHERE
					(SELECT COUNT(".$this->options['structure']['id'].") FROM ".$this->options['structure_table']." WHERE ".$this->options['structure']['id']." = ".$this->options['data_table'].".".$this->options['data2structure'].") = 0
			");
		}
		return true;
	}

	public function res($data = array()) {
		if(!$this->db->query("TRUNCATE TABLE ".$this->options['structure_table'])) { return false; }
		if(!$this->db->query("TRUNCATE TABLE ".$this->options['data_table'])) { return false; }
		$sql = "INSERT INTO ".$this->options['structure_table']." (".implode(",", $this->options['structure']).") VALUES (?".str_repeat(',?', count($this->options['structure']) - 1).")";
		$par = array();
		foreach($this->options['structure'] as $k => $v) {
			switch($k) {
				case 'id':
					$par[] = null;
					break;
				case 'left':
					$par[] = 1;
					break;
				case 'right':
					$par[] = 2;
					break;
				case 'level':
					$par[] = 0;
					break;
				case 'parent_id':
					$par[] = 0;
					break;
				case 'position':
					$par[] = 0;
					break;
				default:
					$par[] = null;
			}
		}
		if(!$this->db->query($sql, $par)) { return false; }
		$id = $this->db->insert_id();
		foreach($this->options['structure'] as $k => $v) {
			if(!isset($data[$k])) { $data[$k] = null; }
		}
		return $this->rn($id, $data);
	}

	public function dump() {
		$nodes = $this->db->get("
			SELECT
				s.".implode(", s.", $this->options['structure']).",
				d.".implode(", d.", $this->options['data'])."
			FROM
				".$this->options['structure_table']." s,
				".$this->options['data_table']." d
			WHERE
				s.".$this->options['structure']['id']." = d.".$this->options['data2structure']."
			ORDER BY ".$this->options['structure']["left"]
		);
		echo "\n\n";
		foreach($nodes as $node) {
			echo str_repeat(" ",(int)$node[$this->options['structure']["level"]] * 2);
			echo $node[$this->options['structure']["id"]]." ".$node["nm"]." (".$node[$this->options['structure']["left"]].",".$node[$this->options['structure']["right"]].",".$node[$this->options['structure']["level"]].",".$node[$this->options['structure']["parent_id"]].",".$node[$this->options['structure']["position"]].")" . "\n";
		}
		echo str_repeat("-",40);
		echo "\n\n";
	}
}