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
  http://127.0.0.1:5000/api/quiz_results?grp=INFSCI2470Spring2023
  ```

- To get reading activities as a CSV file, go to:
  ```
  http://127.0.0.1:5000/api/reading_activities?grp=INFSCI2470Spring2023
  ```

Replace `INFSCI2470Spring2023` with the appropriate group identifier for your data.

## Report Generator Setup and Usage

### Setup

1. Obtain a TSV or CSV file containing the quiz data you want to process.

### Usage

1. Update the `file_path` variable in the `report.py` script with the path to your TSV or CSV file containing quiz data.
2. Set the `delimiter` variable in the `report.py` script to `'\t'` for TSV files or `','` for CSV files.
3. Run the script:

```bash
python report.py
```

The script will read the data, process it, print the results, and save the results to an Excel file named `results.xlsx`. The Excel file will have two sheets:

- **User_DocSrc**: Summary of users' correct answer per document source.
- **User_DocNo**: Users's correct answers for each document, along with the maximum possible correct answers per document in the second row.
