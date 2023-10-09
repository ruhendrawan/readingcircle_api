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


function process_quiz_results($data) {
    $user_docno_correct = array();
    $user_docsrc_correct = array();
    $docno_question_counts = array();
    $docsrc_question_counts = array();

    foreach ($data as $row) {
        $usr = $row['usr'];
        $docno = $row['docno'];
        $docsrc = $row['docsrc'];
        $correct = (int)$row['correct'];

        if (!isset($user_docno_correct[$usr])) {
            $user_docno_correct[$usr] = array();
        }
        if (!isset($user_docno_correct[$usr][$docno])) {
            $user_docno_correct[$usr][$docno] = 0;
        }
        $user_docno_correct[$usr][$docno] += $correct;

        if (!isset($user_docsrc_correct[$usr])) {
            $user_docsrc_correct[$usr] = array();
        }
        if (!isset($user_docsrc_correct[$usr][$docsrc])) {
            $user_docsrc_correct[$usr][$docsrc] = 0;
        }
        $user_docsrc_correct[$usr][$docsrc] += $correct;

        if (!isset($docno_question_counts[$docno])) {
            $docno_question_counts[$docno] = 0;
        }
        $docno_question_counts[$docno] = max($docno_question_counts[$docno], $user_docno_correct[$usr][$docno]);

        if (!isset($docsrc_question_counts[$docsrc])) {
            $docsrc_question_counts[$docsrc] = 0;
        }
        $docsrc_question_counts[$docsrc] = max($docsrc_question_counts[$docsrc], $user_docsrc_correct[$usr][$docsrc]);
    }

    return array($user_docno_correct, $user_docsrc_correct, $docno_question_counts, $docsrc_question_counts);
}




function write_quiz_results_to_csv($user_docno_correct, $user_docsrc_correct, $docno_question_counts, $docsrc_question_counts) {
    $csv_output = fopen('php://temp', 'r+');

    // Write User_DocNo data
    $header = array_merge(array("User"), array_keys($docno_question_counts));
    fputcsv($csv_output, $header);
    $max_row = array_merge(array("Max"), array_values($docno_question_counts));
    fputcsv($csv_output, $max_row);
    foreach ($user_docno_correct as $usr => $docno_correct) {
        $user_row = array($usr);
        foreach ($docno_question_counts as $docno => $count) {
            $user_row[] = isset($docno_correct[$docno]) ? $docno_correct[$docno] : 0;
        }
        fputcsv($csv_output, $user_row);
    }

    // Write User_DocSrc data
    fputcsv($csv_output, array());  // Empty row for separation
    $header = array_merge(array("User"), array_keys($docsrc_question_counts));
    fputcsv($csv_output, $header);
    $max_row = array_merge(array("Max"), array_values($docsrc_question_counts));
    fputcsv($csv_output, $max_row);
    foreach ($user_docsrc_correct as $usr => $docsrc_correct) {
        $user_row = array($usr);
        foreach ($docsrc_question_counts as $docsrc => $count) {
            $user_row[] = isset($docsrc_correct[$docsrc]) ? $docsrc_correct[$docsrc] : 0;
        }
        fputcsv($csv_output, $user_row);
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
} elseif ($_GET['api'] === 'sum_quiz') {
    list($header, $data) = get_quiz_results($grp);
    list($user_docno_correct, $user_docsrc_correct, $docno_question_counts, $docsrc_question_counts) = process_quiz_results($data);
    $csv_data = write_quiz_results_to_csv($user_docno_correct, $user_docsrc_correct, $docno_question_counts, $docsrc_question_counts);
    header('Content-Type: text/csv');
    header('Content-Disposition: attachment; filename="quiz_summary.csv"');
    echo $csv_data;
}

} else {
    echo json_encode(array("error" => "Please provide a 'grp' parameter."));
}

?>
