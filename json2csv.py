#!/usr/bin/env python3

import json
import csv

def json_lines_to_csv(input_file, output_file):
    # Open input file containing newline-separated JSON lines
    with open(input_file, 'r') as f:
        # Read all lines from the file
        json_lines = f.readlines()

    if not json_lines:
        print("No data found in the input file.")
        return

    # Parse the first JSON line to get column headers
    first_line_data = json.loads(json_lines[0])
    column_headers = list(first_line_data.keys())

    # Open output CSV file for writing
    with open(output_file, 'w', newline='') as csv_file:
        writer = csv.DictWriter(csv_file, fieldnames=column_headers)

        # Write column headers to CSV file
        writer.writeheader()

        # Process each JSON line and write to CSV
        for line in json_lines:
            json_data = json.loads(line)
            # Extract values in the same order as column_headers
            row_data = {key: json_data.get(key, '') for key in column_headers}
            writer.writerow(row_data)

    print(f"CSV file '{output_file}' has been generated successfully.")


# first command line argument is the input file
# second command line argument is the output file
if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: python json2csv.py <input_file> <output_file>")
        sys.exit(1)


    # if no second argument, have hte output csv be the same name as the input json file
    output_file = sys.argv[1].split('.')[0] + ".csv" if len(sys.argv) < 3 else sys.argv[2]
    json_lines_to_csv(sys.argv[1], output_file)

