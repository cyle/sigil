<?php

$dbhost = 'http://localhost:8777';

// save to db function
function databaseCall($path = '/', $wut = null, $getheaders = true) {
	global $dbhost;
	if (is_array($wut)) {
		$body = json_encode($wut);
	} else if ($wut != null) {
		$body = trim($wut);
	}
	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $dbhost . $path);
	curl_setopt($ch, CURLOPT_HEADER, $getheaders);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	curl_setopt($ch, CURLOPT_TIMEOUT, 1);
	if (isset($body) && $body != '') {
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');
		curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
	} else {
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'GET');
	}
	$raw_result = curl_exec($ch);
	curl_close($ch);
	if ($raw_result == '') {
		return false;
	}
	$json_result = json_decode($raw_result, true);
	if (is_array($json_result)) {
		return $json_result;
	} else {
		return $raw_result;
	}
}

if (databaseCall('/') == '') {
	die('database not accessible, make sure it\'s running');
}

// ok -- get stuff
$nodes = databaseCall('/nodes', null, false);
if (!is_array($nodes)) {
	die('no nodes to show');
}
$conns = databaseCall('/connections', null, false);

?>
<!--
<?php
print_r($nodes);
print_r($conns);
?>
-->
<?php

function getNodeCoords($node_id) {
	global $nodes;
	$coords = array('X' => 0, 'Y' => 0);
	foreach ($nodes as $node) {
		if ($node['Id'] * 1 == $node_id) {
			$coords['X'] = $node['X'] * 1;
			$coords['Y'] = $node['Y'] * 1;
		}
	}
	return $coords;
}

$grid_scale = 25; // this many pixels for every grid point
$grid_width = 0; // will be set in a minute
$grid_height = 0; // will be set in a minute

// get max node X and add 1 to make the max width
// get max node Y and add 1 to make the max height
foreach ($nodes as $node) {
	if ($node['X'] * 1 >= $grid_width) {
		$grid_width = $node['X'] + 2;
	}
	if ($node['Y'] * 1 >= $grid_height) {
		$grid_height = $node['Y'] + 2;
	}
}

$actual_grid_width = $grid_width * $grid_scale;
$actual_grid_height = $grid_height * $grid_scale;

?><!doctype html>
<html>
<head>
<title>Visualize the Graph</title>
</head>
<body>

<h1>SIGIL Visualization</h1>

<div>
<canvas id="graph" style="background-color:#fff;" width="<?php echo $actual_grid_width; ?>px" height="<?php echo $actual_grid_height; ?>px"></canvas>
</div>

<script type="text/javascript">
// set up the canvas
var a_canvas = document.getElementById('graph');
var a = a_canvas.getContext('2d');
a.font = "12px Courier New";

// draw gridlines
a.strokeStyle = "#ddd";
a.lineWidth = 1;
for (var x = 0.5; x < <?php echo $actual_grid_width; ?>; x += <?php echo $grid_scale; ?>) {
	a.beginPath();
	a.moveTo(x, 0);
	a.lineTo(x, <?php echo $actual_grid_height; ?>);
	a.stroke();
	a.closePath();
}
for (var y = 0.5; y < <?php echo $actual_grid_height; ?>; y += <?php echo $grid_scale; ?>) {
	a.beginPath();
	a.moveTo(0, y);
	a.lineTo(<?php echo $actual_grid_width; ?>, y);
	a.stroke();
	a.closePath();
}

// draw connections
a.strokeStyle = "#900";
a.lineWidth = 2;
<?php
if (is_array($conns)) {
	foreach ($conns as $conn) {
		$source_coords = getNodeCoords($conn['Source']);
		$target_coords = getNodeCoords($conn['Target']);
		echo 'a.beginPath();'."\n";
		echo 'a.moveTo('.($source_coords['X'] * $grid_scale).', '.($source_coords['Y'] * $grid_scale).');'."\n"; // source
		echo 'a.lineTo('.($target_coords['X'] * $grid_scale).', '.($target_coords['Y'] * $grid_scale).');'."\n"; // target
		echo 'a.stroke();'."\n";
		echo 'a.closePath();'."\n";
	}
}
?>

// draw nodes
<?php
foreach ($nodes as $node) {
	echo 'a.beginPath();'."\n";
	echo 'a.arc('.($node['X'] * $grid_scale).', '.($node['Y'] * $grid_scale).', 5, 0, 2 * Math.PI, false);'."\n";
	echo 'a.fillStyle = "#aaa";'."\n"; // fill of the node itself
	echo 'a.fill();'."\n";
	echo 'a.closePath();'."\n";
	echo 'a.fillStyle = "#ccc";'."\n";
	echo 'a.textAlign = "center";'."\n";
	echo 'a.fillText("'.$node['Name'].'", '.($node['X'] * $grid_scale).', '.($node['Y'] * $grid_scale).');'."\n";
	echo 'a.fillStyle = "#000";'."\n";
	echo 'a.fillText("'.$node['Name'].'", '.($node['X'] * $grid_scale - 1).', '.($node['Y'] * $grid_scale - 1).');'."\n";
}
?>

// that's all, folks
</script>
</body>
</html>