<?php
error_reporting(E_ALL);

// Database configuration
define('HOST', 'localhost');
define('DATABASE', 'readingcircle_dev');
define('USER', 'root');
define('PASSWORD', '');

// Connect to the database
function connect_db() {
    $conn = new mysqli(HOST, USER, PASSWORD, DATABASE);
    if ($conn->connect_error) {
        die("Connection failed: " . $conn->connect_error);
    }
    return $conn;
}

// Get quiz results
function get_quiz_results($grp) {
    $conn = connect_db();
    $stmt = $conn->prepare("
        SELECT * FROM submitted_answers s
        INNER JOIN questions q ON s.idquestions = q.idquestions
        INNER JOIN document d ON d.docid = q.docid
        WHERE s.grp = ?
        ORDER BY s.time DESC;
    ");
    $stmt->bind_param("s", $grp);
    $stmt->execute();

    $meta = $stmt->result_metadata();
    $header = array();
    while ($field = $meta->fetch_field()) {
        $header[] = $field->name;
        $params[] = &$row[$field->name];
    }

    call_user_func_array(array($stmt, 'bind_result'), $params);

    $data = array();
    while ($stmt->fetch()) {
        $data[] = array_map(function($value) { return $value; }, $row);
    }

    $stmt->close();
    $conn->close();
    return array($header, $data);
}

// Get reading activities
function get_reading_activities($grp) {
    $conn = connect_db();
    $stmt = $conn->prepare("
        SELECT * FROM tracking
        WHERE grp = ?
    ");
    $stmt->bind_param("s", $grp);
    $stmt->execute();

    $meta = $stmt->result_metadata();
    $header = array();
    while ($field = $meta->fetch_field()) {
        $header[] = $field->name;
        $params[] = &$row[$field->name];
    }

    call_user_func_array(array($stmt, 'bind_result'), $params);

    $data = array();
    while ($stmt->fetch()) {
        $data[] = array_map(function($value) { return $value; }, $row);
    }

    $stmt->close();
    $conn->close();
    return array($header, $data);
}

// Convert data to CSV
function to_csv($header, $results) {
    $csv_output = fopen('php://temp', 'r+');
    fputcsv($csv_output, $header);
    foreach ($results as $row) {
        fputcsv($csv_output, $row);
    }
    rewind($csv_output);
    $csv_data = stream_get_contents($csv_output);
    fclose($csv_output);
    return $csv_data;
}

// Main logic to process the data and generate CSV
if (isset($_GET['grp'])) {
    $grp = $_GET['grp'];
    if ($_GET['api'] === 'raw_quiz') {
        list($header, $results) = get_quiz_results($grp);
        $csv_data = to_csv($header, $results);
        header('Content-Type: text/csv');
        header('Content-Disposition: attachment; filename="quiz_results.csv"');
        echo $csv_data;
    } elseif ($_GET['api'] == 'raw_reading') {
        list($header, $results) = get_reading_activities($grp);
        $csv_data = to_csv($header, $results);
        header('Content-Type: text/csv');
        header('Content-Disposition: attachment; filename="reading_activities.csv"');
        echo $csv_data;
    }
} else {
    echo json_encode(array("error" => "Please provide a 'grp' parameter."));
}

?>
