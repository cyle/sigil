<?php

$dbhost = 'http://localhost:8777';

// save to db function
function databaseCall($path = '/', $wut = null) {
	global $dbhost;
	if (is_array($wut)) {
		$body = json_encode($wut);
	} else if ($wut != null) {
		$body = trim($wut);
	}
	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $dbhost . $path);
	curl_setopt($ch, CURLOPT_HEADER, true);
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
	$json_result = json_decode($raw_result, true);
	if (is_array($json_result)) {
		return $json_result;
	} else {
		return $raw_result;
	}
}

// make the test node map with grid

//$thedata = array();
//$thedata['Name'] = 'The Graph/Spatial Database';

$nodes = array();
$connections = array();


// build the demo nodes

// nodes need these fields: Id, Name, X, Y, Z, ExtraJSONBytes, ExtraJSON
$nodes[] = array( 'Id' => 1, 'Name' => 'Node 1', 'X' => 2, 'Y' => 6);
$nodes[] = array( 'Id' => 2, 'Name' => 'Node 2', 'X' => 3, 'Y' => 3);
$nodes[] = array( 'Id' => 3, 'Name' => 'Node 3', 'X' => 7, 'Y' => 3);
$nodes[] = array( 'Id' => 4, 'Name' => 'Node 4', 'X' => 9, 'Y' => 2);
$nodes[] = array( 'Id' => 5, 'Name' => 'Node 5', 'X' => 13, 'Y' => 5);
$nodes[] = array( 'Id' => 6, 'Name' => 'Node 6', 'X' => 9, 'Y' => 6);
$nodes[] = array( 'Id' => 7, 'Name' => 'Node 7', 'X' => 11, 'Y' => 8);
$nodes[] = array( 'Id' => 8, 'Name' => 'Node 8', 'X' => 12, 'Y' => 11);
$nodes[] = array( 'Id' => 9, 'Name' => 'Node 9', 'X' => 8, 'Y' => 10);
$nodes[] = array( 'Id' => 10, 'Name' => 'Node 10', 'X' => 5, 'Y' => 8);
$nodes[] = array( 'Id' => 11, 'Name' => 'Node 11', 'X' => 5, 'Y' => 11);
$nodes[] = array( 'Id' => 12, 'Name' => 'Node 12', 'X' => 10, 'Y' => 13);
$nodes[] = array( 'Id' => 13, 'Name' => 'Node 13', 'X' => 7, 'Y' => 13);
$nodes[] = array( 'Id' => 14, 'Name' => 'Node 14', 'X' => 4, 'Y' => 13);

// give them all the extra zero'd-out stuff
for ($i = 0; $i < count($nodes); $i++) {
	$nodes[$i]['Z'] = 0.0;
	$nodes[$i]['ExtraJSON'] = null;
}


// build the demo connections

// connections need these fields: Id, Name, Source, Target, Distance, DistanceMultiplier
$connections[] = array( 'Id' => 1, 'Name' => 'Node 1 to 2', 'Source' => 1, 'Target' => 2 );
$connections[] = array( 'Id' => 2, 'Name' => 'Node 2 to 3', 'Source' => 2, 'Target' => 3 );
$connections[] = array( 'Id' => 3, 'Name' => 'Node 3 to 4', 'Source' => 3, 'Target' => 4 );
$connections[] = array( 'Id' => 4, 'Name' => 'Node 4 to 5', 'Source' => 4, 'Target' => 5 );
$connections[] = array( 'Id' => 5, 'Name' => 'Node 3 to 6', 'Source' => 3, 'Target' => 6 );
$connections[] = array( 'Id' => 6, 'Name' => 'Node 5 to 6', 'Source' => 5, 'Target' => 6 );
$connections[] = array( 'Id' => 7, 'Name' => 'Node 6 to 7', 'Source' => 6, 'Target' => 7 );
$connections[] = array( 'Id' => 8, 'Name' => 'Node 7 to 8', 'Source' => 7, 'Target' => 8 );
$connections[] = array( 'Id' => 9, 'Name' => 'Node 7 to 9', 'Source' => 7, 'Target' => 9 );
$connections[] = array( 'Id' => 10, 'Name' => 'Node 8 to 12', 'Source' => 8, 'Target' => 12 );
$connections[] = array( 'Id' => 11, 'Name' => 'Node 12 to 13', 'Source' => 12, 'Target' => 13 );
$connections[] = array( 'Id' => 12, 'Name' => 'Node 13 to 14', 'Source' => 13, 'Target' => 14 );
$connections[] = array( 'Id' => 13, 'Name' => 'Node 13 to 11', 'Source' => 13, 'Target' => 11 );
$connections[] = array( 'Id' => 14, 'Name' => 'Node 11 to 14', 'Source' => 11, 'Target' => 14 );
$connections[] = array( 'Id' => 15, 'Name' => 'Node 11 to 10', 'Source' => 11, 'Target' => 10 );
$connections[] = array( 'Id' => 16, 'Name' => 'Node 10 to 9', 'Source' => 10, 'Target' => 9 );

// calculate distances between nodes for the connections
for ($i = 0; $i < count($connections); $i++) {
	
	$distance = 0;
	$this_conn = $connections[$i];
	
	$x1 = 0;
	$y1 = 0;
	
	$x2 = 0;
	$y2 = 0;
	
	// get source node's coordinates
	foreach ($nodes as $node) {
		if ($node['Id'] == $this_conn['Source']) {
			$x1 = $node['X'];
			$y1 = $node['Y'];
		}
	}
	
	// get target node's coordinates
	foreach ($nodes as $node) {
		if ($node['Id'] == $this_conn['Target']) {
			$x2 = $node['X'];
			$y2 = $node['Y'];
		}
	}
	
	// get distance
	$xd = $x2 - $x1;
	$yd = $y2 - $y1;
	$distance = sqrt( ($xd * $xd) + ($yd * $yd) );
	
	$connections[$i]['Distance'] = $distance;
}

// give them all the extra zero'd-out stuff
for ($i = 0; $i < count($connections); $i++) {
	$connections[$i]['DistanceMultiplier'] = null;
	$connections[$i]['ExtraJSON'] = null;
}

//$thedata['Connections'] = $nodes;
//$thedata['Nodes'] = $connections;

// save them to the database
foreach ($nodes as $node) {
	$dbresult = databaseCall('/node', $node);
	if (substr($dbresult, 0, 20) != 'HTTP/1.1 201 Created') {
		die('error: '.$dbresult);
	}
}

foreach ($connections as $connection) {
	$dbresult = databaseCall('/connection', $connection);
	if (substr($dbresult, 0, 20) != 'HTTP/1.1 201 Created') {
		die('error: '.$dbresult);
	}
}

echo 'okay'."\n";

?>