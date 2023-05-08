import argparse
import csv

from lib_report import process_quiz_results, write_quiz_results_to_excel


def read_data(file_path, delimiter):
    data = []
    with open(file_path, newline='', encoding='utf-8', errors='replace') as csvfile:
        reader = csv.DictReader(csvfile, delimiter=delimiter)
        for row in reader:
            data.append(row)
    return data


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Process quiz data and generate a report.")
    parser.add_argument("input_file", help="Path to the input TSV or CSV file containing quiz data.")
    parser.add_argument("output_file", help="Path where the generated Excel file will be saved.")
    parser.add_argument("-d", "--delimiter", default=',', help="Delimiter for the input file. Use '\\t' for TSV files and ',' for CSV files. Default: '\\t'")

    args = parser.parse_args()

    if args.delimiter == '\\t':
        data = read_data(args.input_file, '\t')
    else:
        data = read_data(args.input_file, args.delimiter)
    user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts = process_quiz_results(data)

    write_quiz_results_to_excel(user_docno_correct, user_docsrc_correct, docno_question_counts, docsrc_question_counts, args.output_file)
    print(f"Results saved to {args.output_file}")
