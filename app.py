# pip install flask
# pip install mysql-connector-python

# http://127.0.0.1:5000/api/quiz_results?grp=INFSCI2470Spring2023
# http://127.0.0.1:5000/api/reading_activities?grp=INFSCI2470Spring2023


import csv
from flask import Flask, jsonify, request, Response
import mysql.connector
from flask import send_file
from lib_report import process_quiz_results, write_quiz_results_to_excel
import io

HOST = 'localhost'
DATABASE = 'readingcircle_dev'
USER = 'root'
PASSWORD = ''

app = Flask(__name__)


def get_quiz_results(grp):
    cnx = mysql.connector.connect(user=USER, password=PASSWORD, host=HOST, database=DATABASE)
    cursor = cnx.cursor()

    query = '''
        SELECT * FROM submitted_answers s
        INNER JOIN questions q ON s.idquestions = q.idquestions
        INNER JOIN document d ON d.docid = q.docid
        INNER JOIN doc_lecture dl on dl.docno = d.docno
        WHERE s.grp = %s
        ORDER BY s.time DESC;
    '''
    cursor.execute(query, (grp,))
    results = cursor.fetchall()

    # Get the header dynamically from the cursor description
    header = [column[0] for column in cursor.description]

    cursor.close()
    cnx.close()

    return header, results


def get_reading_activities(grp):
    cnx = mysql.connector.connect(user=USER, password=PASSWORD, host=HOST, database=DATABASE)
    cursor = cnx.cursor()

    query = '''
SELECT * FROM tracking  
WHERE s.grp = %s;
    '''
    cursor.execute(query, (grp,))
    results = cursor.fetchall()

    # Get the header dynamically from the cursor description
    header = [column[0] for column in cursor.description]

    cursor.close()
    cnx.close()

    return header, results


@app.route('/api/raw_quiz_results', methods=['GET'])
def raw_quiz_results():
    grp = request.args.get('grp')
    if grp:
        header, results = get_quiz_results(grp)
        csv_output = to_csv(header, results)

        return Response(
            csv_output,
            mimetype="text/csv",
            headers={"Content-Disposition": "attachment;filename=quiz_results.csv"},
        )
    else:
        return jsonify({"error": "Please provide a 'grp' parameter."}), 400


@app.route('/api/raw_reading_activities', methods=['GET'])
def raw_reading_activities():
    grp = request.args.get('grp')
    if grp:
        header, results = get_reading_activities(grp)
        csv_output = to_csv(header, results)

        return Response(
            csv_output,
            mimetype="text/csv",
            headers={"Content-Disposition": "attachment;filename=reading_activities.csv"},
        )
    else:
        return jsonify({"error": "Please provide a 'grp' parameter."}), 400


@app.route('/api/xls_quiz_results', methods=['GET'])
def xls_quiz_results():
    grp = request.args.get('grp')
    if grp:
        header, results = get_quiz_results(grp)

        data = [dict(zip(header, row)) for row in results]

        user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts = \
            process_quiz_results(data)

        output_file = io.BytesIO()
        write_quiz_results_to_excel(user_docno_correct, user_docsrc_correct,
                                    docno_question_counts, docsrc_question_counts,
                                    output_file)
        output_file.seek(0)

        return send_file(
            output_file,
            as_attachment=True,
            download_name='ISD_rc_results.xlsx',
            mimetype='application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        )
    else:
        return jsonify({"error": "Please provide a 'grp' parameter."}), 400


# Create a CSV file-like object in memory
def to_csv(header, results):
    csv_data = [header] + results
    csv_file = csv.StringIO()
    writer = csv.writer(csv_file)
    writer.writerows(csv_data)
    csv_output = csv_file.getvalue()
    return csv_output


if __name__ == '__main__':
    app.run(debug=True)

