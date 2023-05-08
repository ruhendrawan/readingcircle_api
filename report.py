# pip install openpyxl

import csv
from collections import defaultdict
from openpyxl import Workbook


def read_data(file_path, delimiter):
    data = []
    with open(file_path, newline='', encoding='utf-8', errors='replace') as csvfile:
        reader = csv.DictReader(csvfile, delimiter=delimiter)
        for row in reader:
            data.append(row)
    return data



def process_data(data):
    user_docno_correct = defaultdict(lambda: defaultdict(int))
    user_docsrc_correct = defaultdict(lambda: defaultdict(int))
    docno_question_counts = defaultdict(int)
    docsrc_question_counts = defaultdict(int)

    for row in data:
        user_docno_correct[row['usr']][row['docno']] += int(row['correct'])
        user_docsrc_correct[row['usr']][row['docsrc']] += int(row['correct'])

        docno_question_counts[row['docno']] = max(docno_question_counts[row['docno']], user_docno_correct[row['usr']][row['docno']])
        docsrc_question_counts[row['docsrc']] = max(docsrc_question_counts[row['docsrc']], user_docsrc_correct[row['usr']][row['docsrc']])

    return user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts

def write_to_excel(user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts, output_file):
    wb = Workbook()

    ws1 = wb.active
    ws1.title = "User_DocNo"
    header = ["User"] + [f"{docno}" for docno in sorted(docno_question_counts)]
    ws1.append(header)
    max_row = ["Max"] + [docno_question_counts[docno] for docno in sorted(docno_question_counts)]
    ws1.append(max_row)

    for usr, docno_correct in user_docno_correct.items():
        user_row = [usr] + [docno_correct[docno] for docno in sorted(docno_question_counts)]
        ws1.append(user_row)

    ws2 = wb.create_sheet("User_DocSrc", -1)
    header = ["User"] + [f"{docno}" for docno in sorted(docsrc_question_counts)]
    ws2.append(header)
    max_row = ["Max"] + [docsrc_question_counts[docno] for docno in sorted(docsrc_question_counts)]
    ws2.append(max_row)
    # Add formula to compute the sum of columns B, C, and D
    ws2.cell(row=2, column=len(header) + 1,
             value=f"=SUM(B{2}:D{2})")

    row_number = 3
    for usr, docno_correct in user_docsrc_correct.items():
        user_row = [usr] + [docno_correct[docno] for docno in sorted(docsrc_question_counts)]
        ws2.append(user_row)

        # Add formula to compute the sum of columns B, C, and D
        ws2.cell(row=row_number, column=len(user_row) + 1,
                 value=f"=SUM(B{row_number}:D{row_number})")

        # Add formula to compute column E divided by $E$2
        ws2.cell(row=row_number, column=len(user_row) + 2,
                 value=f"=ROUND(E{row_number}/$E$2*100, 2)")

        row_number += 1

    wb.save(output_file)

if __name__ == "__main__":
    root_path = "/Users/rully/Downloads/"
    file_path = root_path + "ISD_quizzes_count.tsv"

    delimiter = '\t'  # Change this to ',' for CSV files

    data = read_data(file_path, delimiter)
    user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts = process_data(data)

    output_file = root_path + "ISD_rc_results.xlsx"
    write_to_excel(user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts, output_file)
    print(f"Results saved to {output_file}")
