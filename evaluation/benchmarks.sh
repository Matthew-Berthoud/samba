csv_file="benchmark-results.csv"
txt_file="benchmark-results.txt"

cd ../
go test -bench=. > "evaluation/$txt_file"
cd evaluation
echo "Raw output written to $txt_file"

echo "method,ms_per_op" > "$csv_file"
cat benchmark-results.txt | awk '/ns\/op/ { 
    gsub(/^Benchmark/, "", $1);    # Remove "Benchmark" from the start
    gsub(/-8$/, "", $1);           # Remove "-8" from the end
    printf "%s,%.3f\n", $1, $3/1000000 
}' >> "$csv_file"
echo "CSV output written to $csv_file"
