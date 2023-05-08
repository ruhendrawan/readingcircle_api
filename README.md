# readingcircle_api
# Quiz Data API and Report Generator

This project consists of two main components: 
- an API to retrieve quiz results and reading activities as CSV files, and 
- a report generator that processes quiz data and generates an Excel file.

## Dependencies

- Python 3.6+
- openpyxl
- Flask
- mysql-connector-python

To install the required packages, run:

```bash
pip install openpyxl flask mysql-connector-python
```

## API Setup and Usage

### Setup

1. Update the `HOST`, `DATABASE`, `USER`, and `PASSWORD` variables in the `api.py` script with your MySQL database connection information.
2. Make sure your database contains the necessary tables and data.

### Usage

1. Run the API server:

```bash
python api.py
```

2. Access the following endpoints using your browser or a tool like curl or Postman:

- To get quiz results as a CSV file, go to:
  ```
  http://127.0.0.1:5000/api/raw_quiz_results?grp=INFSCI2470Spring2023
  ```

- To get reading activities as a CSV file, go to:
  ```
  http://127.0.0.1:5000/api/raw_reading_activities?grp=INFSCI2470Spring2023
  ```

- To generate an Excel report for quiz results, go to:
  ```
  http://127.0.0.1:5000/api/xls_quiz_results?grp=INFSCI2470Spring2023
  ```

Replace `INFSCI2470Spring2023` with the appropriate group identifier for your data.

## Report Generator Setup and Usage

### Setup

Obtain a TSV or CSV file containing the quiz data you want to process.

### Usage

Run the script:

```bash
python cli-report.py path/to/input_file.tsv path/to/output_file.xlsx
```

E.g.
```bash
python3 cli-report.py /Users/rully/Downloads/quiz_results.csv  /Users/rully/Downloads/ISD_rc_results.xlsx  
```

You can also provide a delimiter using the -d or --delimiter option:

```bash
python cli-report.py -d '\t' path/to/input_file.tsv path/to/output_file.xlsx
```

The script will read the data, process it, print the results, and save the results to an Excel file named `results.xlsx`. The Excel file will have two sheets:

- **User_DocSrc**: Summary of users' correct answer per document source.
- **User_DocNo**: Users's correct answers for each document, along with the maximum possible correct answers per document in the second row.
